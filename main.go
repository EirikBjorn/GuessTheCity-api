package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type City struct {
	ID int `json:"rank"`
	Name string `json:"name"`
	Country string `json:"country"`
	Population int `json:"population"`
}

type Cities struct {
	Collection []City
}

func main() {
	router := gin.Default()
	router.GET("/cities", getCities)
	router.Run("localhost:8080")
}

func getDataFromJson() []City {

	data, err := ioutil.ReadFile("./data.json")
	if err != nil {
		fmt.Println(err)
	}

	cities := make([]City,0)
  json.Unmarshal(data, &cities)

	return cities
}

func getCities(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getDataFromJson())
}