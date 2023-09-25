package app

import (
	"log"
	"net/http"
	"space/internal/app/ds"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (a *Application) StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.Static("/images", "./resources/images")
	r.Static("/fonts", "./resources/fonts")
	r.Static("/data", "./resources/data")
	r.Static("/css", "./resources/css")

	r.GET("/", func(c *gin.Context) {
		planets, err := a.repo.GetActivePlanets()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		data := gin.H{
			"planets": planets,
		}
		c.HTML(http.StatusOK, "mainPage.tmpl", data)
	})

	r.GET("/planets/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		planet, err := a.repo.GetActivePlanetById(id)
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "planet.tmpl", planet)
	})

	r.GET("/search", func(c *gin.Context) {
		planets, err := a.repo.GetActivePlanets()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		planetSlice := *planets
		searchQuery := c.DefaultQuery("q", "")
		var foundPlanets []ds.Planet
		for _, planet := range planetSlice {
			if strings.HasPrefix(strings.ToLower(planet.Name), strings.ToLower(searchQuery)) {
				foundPlanets = append(foundPlanets, planet)
			}
		}
		data := gin.H{
			"planets": foundPlanets,
		}
		c.HTML(http.StatusOK, "mainPage.tmpl", data)
	})

	r.POST("/delete", func(c *gin.Context) {
		id, err := strconv.Atoi(c.DefaultQuery("q", ""))
		log.Print(c.DefaultQuery("q", ""))
		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = a.repo.DeactivatePlanetByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		planets, err := a.repo.GetActivePlanets()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		data := gin.H{
			"planets": planets,
		}
		c.HTML(http.StatusOK, "mainPage.tmpl", data)
	})

	r.Run()
	// go run cmd/space/main.go

	log.Println("Server down")
}
