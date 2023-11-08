package app

import (
	"log"
	"space/internal/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (a *Application) StartServer() {
	log.Println("Server start up")

	handler := api.NewHandler(a.repo)
	r := gin.Default()
	r.Use(cors.Default())

	PlanetGroup := r.Group("/planet")
	{
		PlanetGroup.GET("/", handler.GetPlanets)             // список всех планет
		PlanetGroup.GET("/:id", handler.GetPlanetById)       // одна планета
		PlanetGroup.DELETE("/:id", handler.DeletePlanetById) // удалить планету по ид
		PlanetGroup.PUT("/:id", handler.ChangePlanetById)    // изменить планету по ид
		PlanetGroup.POST("/", handler.CreatePlanet)          // создать планету
		PlanetGroup.POST("/:id", handler.AddPlanetById)      // добавить планету в созвездие
	}

	StellaGroup := r.Group("/constellation")
	{
		StellaGroup.GET("/", handler.GetConstellations)             // все созвездия
		StellaGroup.GET("/:id", handler.GetConstellationById)       // одно созвездие с планетами
		StellaGroup.PUT("/:id", handler.ChangeConstellationById)    // изменить поля созвездия
		StellaGroup.DELETE("/:id", handler.DeleteConstellationById) // удалить созвездие
		StellaGroup.PUT("/inprogress", handler.DoConstellationInProgress)
		StellaGroup.PUT("/cancel/:id", handler.DoConstelltionCanceledById)
		StellaGroup.PUT("/complete/:id", handler.DoConstelltionCompletedById)
		StellaGroup.DELETE("/remove/:id", handler.RemovePlanetById) // удалить планету из созвездия по ид планеты
	}

	MinioGroup := r.Group("/minio")
	{
		MinioGroup.POST("/:id", handler.AddImage) // добавить картинку для планеты
	}

	r.Run()
	// go run cmd/space/main.go

	log.Println("Server down")
}
