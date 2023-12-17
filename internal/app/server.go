package app

import (
	"log"
	"space/docs"
	"space/internal/api"
	"space/internal/app/middleware"
	"space/internal/app/repository"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func (a *Application) StartServer() {
	log.Println("Server start up")

	handler := api.NewHandler(a.repo, repository.NewRedis())
	r := gin.Default()

	docs.SwaggerInfo.Title = "space"
	docs.SwaggerInfo.Description = "rip course project about space and constellations"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = ""

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler)) // http://localhost:8080/swagger/index.html#/
	r.POST("/auth/login", handler.Login)
	r.POST("/auth/logout", handler.Logout)
	r.POST("/auth/register", handler.Register)

	r.GET("/planet", handler.GetPlanets)        // список всех планет
	r.GET("/planet/:id", handler.GetPlanetById) // одна планета
	mw := middleware.New(handler.RedisRepo)
	authorized := r.Group("/")
	authorized.Use(mw.IsAuth())
	{
		PlanetGroup := authorized.Group("/planet")
		{
			PlanetGroup.DELETE("/:id", handler.DeletePlanetById) // удалить планету по ид
			PlanetGroup.PUT("/:id", handler.ChangePlanetById)    // изменить планету по ид
			PlanetGroup.POST("/", handler.CreatePlanet)          // создать планету
			PlanetGroup.POST("/:id", handler.AddPlanetById)      // добавить планету в созвездие
		}

		StellaGroup := authorized.Group("/constellation")
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

		MinioGroup := authorized.Group("/minio")
		{
			MinioGroup.POST("/:id", handler.AddImage) // добавить картинку для планеты
		}

	}
	r.Run()
	// go run cmd/space/main.go

	log.Println("Server down")
}
