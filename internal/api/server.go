package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	file, err := os.Open("resources/data/planets.json")
	if err != nil {
		fmt.Println("Ошибка при открытии JSON файла:", err)
		return
	}
	defer file.Close()

	var planets []Planet
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&planets); err != nil {
		fmt.Println("Ошибка при декодировании JSON данных:", err)
		return
	}

	r.LoadHTMLGlob("templates/*")

	r.Static("/images", "./resources/images")
	r.Static("/fonts", "./resources/fonts")
	r.Static("/data", "./resources/data")
	r.Static("/css", "./resources/css")

	r.GET("/", func(c *gin.Context) {
		data := gin.H{
			"planets": planets,
		}
		c.HTML(http.StatusOK, "mainPage.tmpl", data)
	})

	r.GET("/planets/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Print(err)
		}
		planet := planets[id-1]
		c.HTML(http.StatusOK, "planet.tmpl", planet)
	})

	r.GET("/search", func(c *gin.Context) {
		searchQuery := c.DefaultQuery("q", "")
		var foundPlanets []Planet
		for _, planet := range planets {
			if strings.HasPrefix(strings.ToLower(planet.Name), strings.ToLower(searchQuery)) {
				foundPlanets = append(foundPlanets, planet)
			}
		}
		data := gin.H{
			"planets": foundPlanets,
		}
		c.HTML(http.StatusOK, "mainPage.tmpl", data)
	})

	r.GET("/show-planets", func(c *gin.Context) {
		var solarSystemPlanets []Planet
		for _, planet := range planets {
			if planet.Id >= 1 && planet.Id <= 9 {
				solarSystemPlanets = append(solarSystemPlanets, planet)
			}
		}
		data := gin.H{
			"planets": solarSystemPlanets,
		}
		c.HTML(http.StatusOK, "mainPage.tmpl", data)
	})

	r.Run()
	// go run cmd/main.go

	log.Println("Server down")
}
