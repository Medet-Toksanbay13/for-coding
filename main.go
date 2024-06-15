package main

import (
	"database/sql"
	"log"
	"net/http"
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

	_, err = db.Exec("INSERT INTO public.forsedb (name, car, buy) VALUES ($1, $2, $3)", newUser.Name, newUser.Car, newUser.Buy)
	if err != nil {
		log.Println("Database insert error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error inserting values", "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"msg": "User added successfully"})
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
	route.Run(":8088")
}
