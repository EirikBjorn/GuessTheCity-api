package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type City struct {
	Name       string `json:"name"`
	Country    string `json:"country"`
	Rank       int    `json:"rank"`
	Population int    `json:"population"`
	Continent  string `json:"continent"`
	IsCapital  bool   `json:"isCapital"`
}

func main() {

	getDataFromJson()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/citiesAll", getCitiesAll)
	router.GET("/citiesLarge", getCitiesLarge)
	router.GET("/citiesPrimate", getCitiesPrimate)
	router.GET("/citiesCapitals", getCitiesCapitals)
	router.GET("/citiesEurope", getCitiesEurope)
	router.GET("/citiesEuropeCapitals", getCitiesEuropeCapitals)
	router.GET("/citiesUS", getCitiesUS)
	router.GET("/citiesNA", getCitiesNA)

	router.GET("/answer/:answer/:rank", func(c *gin.Context) {
		ans, rankStr := c.Param("answer"), c.Param("rank")
		rank, err := strconv.Atoi(rankStr)
		correct := false
		fmt.Print(err)
		for _, i := range getDataFromJson() {
			if strings.EqualFold(ans, i.Name) && rank == i.Rank {
				correct = true
			}
		}
		c.JSON(200, gin.H{
			"message": correct,
		})
	})

	router.Run()
}

// Read json containing cities and return them as an slice
func getDataFromJson() []City {
	data, err := ioutil.ReadFile("data.json")
	if err != nil {
		fmt.Println(err)
	}
	var citiesList []City
	json.Unmarshal(data, &citiesList)
	return citiesList
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
		}
	}
	return finishedSlice
}

func antiCheat(s []City) []City {
	var antiCheatSlice []City
	for _, element := range s {
		element.Name = "No cheating ;)"
		element.Country = "No cheating ;)"
		element.Population = 1337
		element.Continent = "No cheating ;)"
		element.IsCapital = false
		antiCheatSlice = append(antiCheatSlice, element)
	}
	return antiCheatSlice
}

// Return 5 random cities
func getCitiesAll(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, antiCheat(shuffleAndShorten(getDataFromJson())))
}

// Return 5 random cities with a population above 10 million
func getCitiesLarge(c *gin.Context) {
	var finishedSlice []City
	for _, element := range getDataFromJson() {
		if element.Population >= 10000000 {
			finishedSlice = append(finishedSlice, element)
		}
	}
	c.IndentedJSON(http.StatusOK, antiCheat(shuffleAndShorten(finishedSlice)))
}

// Return 5 random primate cities
func getCitiesPrimate(c *gin.Context) {
	var finishedSlice []City
	for _, element := range getDataFromJson() {
		if !containsCountry(finishedSlice, element) {
			finishedSlice = append(finishedSlice, element)
		}
	}
	c.IndentedJSON(http.StatusOK, antiCheat(shuffleAndShorten(finishedSlice)))
}

// Return 5 random capitals
func getCitiesCapitals(c *gin.Context) {
	var finishedSlice []City
	for _, element := range getDataFromJson() {
		if element.IsCapital {
			finishedSlice = append(finishedSlice, element)
		}
	}
	c.IndentedJSON(http.StatusOK, antiCheat(shuffleAndShorten(finishedSlice)))
}

// Return 5 random european cities (Exclusing russia because there is soooo many cities)
func getCitiesEurope(c *gin.Context) {
	var finishedSlice []City
	for _, element := range getDataFromJson() {
		if element.Continent == "Europe" && element.Country != "Russia" {
			finishedSlice = append(finishedSlice, element)
		}
	}
	c.IndentedJSON(http.StatusOK, antiCheat(shuffleAndShorten(finishedSlice)))
}

// Return 5 random european capitals
func getCitiesEuropeCapitals(c *gin.Context) {
	var finishedSlice []City
	for _, element := range getDataFromJson() {
		if element.Continent == "Europe" && element.IsCapital {
			finishedSlice = append(finishedSlice, element)
		}
	}
	c.IndentedJSON(http.StatusOK, antiCheat(shuffleAndShorten(finishedSlice)))
}

// Return 5 random US Cities
func getCitiesUS(c *gin.Context) {
	var finishedSlice []City
	for _, element := range getDataFromJson() {
		if element.Country == "United States" {
			finishedSlice = append(finishedSlice, element)
		}
	}
	c.IndentedJSON(http.StatusOK, antiCheat(shuffleAndShorten(finishedSlice)))
}

// Return 5 random US Cities
func getCitiesNA(c *gin.Context) {
	var finishedSlice []City
	for _, element := range getDataFromJson() {
		if element.Continent == "North America" {
			finishedSlice = append(finishedSlice, element)
		}
	}
	c.IndentedJSON(http.StatusOK, antiCheat(shuffleAndShorten(finishedSlice)))
}
