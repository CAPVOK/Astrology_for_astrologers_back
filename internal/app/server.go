package app

import (
	"fmt"
	"log"

	"space/docs"
	"space/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title BagTracker RestAPI
// @version 1.0
// @description API server for Space application

// @host http://localhost:8081
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// Run запускает приложение.
func (app *Application) Run() {
	r := gin.Default()
	// Это нужно для автоматического создания папки "docs" в вашем проекте
	docs.SwaggerInfo.Title = "BagTracker RestAPI"
	docs.SwaggerInfo.Description = "API server for BagTracker application"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8081"
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Группа запросов для планет
	PlanetGroup := r.Group("/planet")
	{
		PlanetGroup.GET("/", middleware.Guest(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.GetPlanets)
		PlanetGroup.GET("/:planet_id", middleware.Guest(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.GetPlanetByID)
		PlanetGroup.DELETE("/:planet_id", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.DeletePlanet)
		PlanetGroup.POST("/", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.CreatePlanet)
		PlanetGroup.PUT("/:planet_id", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.UpdatePlanet)
		PlanetGroup.POST("/:planet_id/constellation", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.AddPlanetToConstellation)
		PlanetGroup.DELETE("/:planet_id/constellation", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.RemovePlanetFromConstellation)
		PlanetGroup.POST("/:planet_id/image", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.AddPlanetImage)
	}

	// Группа запросов для созвездий
	ConstellationGroup := r.Group("/constellation").Use(middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository))
	{
		ConstellationGroup.GET("/", app.Handler.GetConstellations)
		ConstellationGroup.GET("/:id", app.Handler.GetConstellationByID)
		ConstellationGroup.DELETE("/:id", app.Handler.DeleteConstellation)
		ConstellationGroup.PUT("/:id/update", app.Handler.UpdateConstellation)
		ConstellationGroup.PUT("/:id/status", app.Handler.UpdateConstellationStatus)
	}

	UserGroup := r.Group("/user")
	{
		UserGroup.POST("/register", app.Handler.Register)
		UserGroup.POST("/login", app.Handler.Login)
		UserGroup.POST("/logout", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.Logout)
	}
	addr := fmt.Sprintf("%s:%d", app.Config.ServiceHost, app.Config.ServicePort)
	r.Run(addr)
	log.Println("Server down")
}
