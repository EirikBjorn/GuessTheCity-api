package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
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
	router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET"},
    AllowHeaders:     []string{"Origin"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge: 12 * time.Hour,
  }))

	router.GET("/citiesAll", getCitiesAll)
	router.GET("/citiesLarge", getCitiesLarge)
	router.GET("/citiesPrimate", getCitiesPrimate)
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

// Generate a random number between 0 and length of slice
func random(in []City) int {
	randomIndex := rand.Intn(len(in))
	return randomIndex
}

// Returns true if a city belonging to same country already exists in the slice (For returning primate cities)
func containsCountry(s []City, e City) bool {
	for _, a := range s {
			if a.Country == e.Country {
				return true
			} 
	}
	return false
}

// Check if city already exists in slice (Prevent duplicate)
func containsCity(s []City, e City) bool {
	for _, a := range s {
			if a.Name == e.Name {
				return true
			} 
	}
	return false
}

// Takes a list of cities and shuffles it and removes all but 5 random cities
func shuffleAndShorten(list []City) []City {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })
	var finishedSlice []City
	for i := 0; i < 5; {
		if !containsCity(finishedSlice, list[i]) {
			finishedSlice = append(finishedSlice, list[i])
			i++
		} else {
			fmt.Print("Duplicate")
		}
	}
	return finishedSlice
}

// Return 5 random cities
func getCitiesAll(c *gin.Context) {
	var finishedSlice []City
	for i := 0; i < 5; i++ {
		finishedSlice = append(finishedSlice, getDataFromJson()[random(getDataFromJson())])
	}
	c.IndentedJSON(http.StatusOK, finishedSlice)
}

// Return 5 random cities with a population above 5 million
func getCitiesLarge(c *gin.Context) {
	var sizedSlice []City
	for _, element := range getDataFromJson() {
    if element.Population >= 5000000 {
			sizedSlice = append(sizedSlice, element)
		}
	}
	c.IndentedJSON(http.StatusOK, shuffleAndShorten(sizedSlice))
}

// Return 5 random primate cities
func getCitiesPrimate(c *gin.Context) {

	var primateSlice []City

	for _, element := range getDataFromJson() {
    if !containsCountry(primateSlice, element) {
			primateSlice = append(primateSlice, element)
		} 
	}
	c.IndentedJSON(http.StatusOK, shuffleAndShorten(primateSlice))
}


