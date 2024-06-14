package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Birthday int    `json:"birthday"`
}

var person = []Person{
	{Name: "Medet", Address: "Almaty", Birthday: 2024},
	{Name: "Someone", Address: "Astana", Birthday: 2022},
}

func main() {
	route := gin.Default()
	route.GET("/testing", startPage)
	route.GET("/testing/:name", findPage)
	route.DELETE("/testing/:name", deletePage)
	route.POST("/testing", postPage)
	route.PUT("/testing/:address", updatePage)
	route.Run(":8085")
}

func startPage(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, person)
}

func postPage(c *gin.Context) {
	var newPerson Person

	err := c.BindJSON(&newPerson)
	if err != nil {
		return
	}
	person = append(person, newPerson)
	c.IndentedJSON(http.StatusCreated, newPerson)
}

func findPage(c *gin.Context) {
	addr := c.Param("name")

	for _, el := range person {
		if el.Name == addr {
			c.IndentedJSON(http.StatusOK, el)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Your input name not eq for element in slice person when finding"})
}

func deletePage(c *gin.Context) {
	name := c.Param("name")

	for i := range person {
		if person[i].Name == name {
			person = append(person[:i], person[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Your input name not eq for element in slice person when delete"})
}

func updatePage(c *gin.Context) {
	address := c.Param("address")

	for i := range person {
		if person[i].Address == address {
			if err := c.BindJSON(&person[i]); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.IndentedJSON(http.StatusOK, person[i])
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"msg": "Your input address not eq for element in slice person when update"})
}
