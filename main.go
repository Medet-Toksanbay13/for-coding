package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type User struct {
	Id   int       `json:"id"`
	Name string    `json:"name"`
	Car  string    `json:"car"`
	Buy  time.Time `json:"buy"`
}

var db *sql.DB

func mainP(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM public.forsedb")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error when selecting from table", "error": err.Error()})
		return

	}
	defer rows.Close()

	var u User
	users := []User{}
	for rows.Next() {
		err = rows.Scan(&u.Id, &u.Name, &u.Car, &u.Buy)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Error when scanning from select", "error": err.Error()})
			return
		}
		users = append(users, u)
	}
	if rows.Err() != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error after scanning from select", "error": rows.Err().Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, users)
}

func createP(c *gin.Context) {
	var newUser User

	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error binding newperson", "error": err.Error()})
		return
	}

	_, err = db.Exec("INSERT INTO public.forsedb (id, name, car, buy) VALUES ($1, $2, $3, $4)", newUser.Id, newUser.Name, newUser.Car, newUser.Buy)
	if err != nil {
		log.Println("Database insert error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error inserting values", "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"msg": "User added successfully"})
}

func deleteP(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error when strconv id in delete", "Error": err.Error()})
		return
	}

	res, err := db.Exec("DELETE FROM public.forsedb WHERE id = $1", idInt)
	if err != nil {
		log.Println("Database delete error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error when deleteing user", "Error": err.Error()})
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Error when rowsAff == 0"})
		return
	}

	log.Println("User deleted successfully:", idInt)
	c.Status(http.StatusNoContent)
}

func updateP(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error when strconv id in delete", "Error": err.Error()})
		return
	}

	var newUser User
	err = c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error binding newperson", "error": err.Error()})
		return
	}

	_, err = db.Exec("UPDATE public.forsedb SET name=$1, car=$2, buy=$3 WHERE id=$4", newUser.Name, newUser.Car, newUser.Buy, idInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "User updated successfully!"})
}

func main() {
	conn := "user=postgres password=meda13 dbname=forsedb sslmode=disable"
	var err error
	db, err = sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err, "Error when connecting with database")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("No database connection:", err)
	}

	defer db.Close()

	route := gin.Default()

	route.GET("/main", mainP)
	route.POST("/main/create", createP)
	route.DELETE("/main/:id", deleteP)
	route.PUT("/main/:id", updateP)
	route.Run(":8088")
}
