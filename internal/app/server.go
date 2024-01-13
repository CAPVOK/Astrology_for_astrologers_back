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
	// Группа запросов для планеты
	PlanetGroup := r.Group("/planet")
	{
		PlanetGroup.GET("/", middleware.Guest(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.GetPlanets)
		PlanetGroup.GET("/:planet_id", middleware.Guest(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.GetPlanetByID)
		PlanetGroup.DELETE("/:planet_id/delete", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.DeletePlanet)
		PlanetGroup.POST("/create", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.CreatePlanet)
		PlanetGroup.PUT("/:planet_id/update", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.UpdatePlanet)
		PlanetGroup.POST("/:planet_id/constellation", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.AddPlanetToConstellation)
		PlanetGroup.DELETE("/:planet_id/constellation/delete", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.RemovePlanetFromConstellation)
		PlanetGroup.POST("/:planet_id/image", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.AddPlanetImage)
	}

	// Группа запросов для созвездия
	ConstellationGroup := r.Group("/constellation").Use(middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository))
	{
		ConstellationGroup.GET("/", app.Handler.GetConstellations)
		ConstellationGroup.GET("/:id", app.Handler.GetConstellationByID)
		ConstellationGroup.DELETE("/:id/delete", app.Handler.DeleteConstellation)
		ConstellationGroup.PUT("/:id/update", app.Handler.UpdateConstellationFlightNumber)
		ConstellationGroup.PUT("/:id/status/user", app.Handler.UpdateConstellationStatusUser)           // Новый маршрут для обновления статуса созвездия пользователем
		ConstellationGroup.PUT("/:id/status/moderator", app.Handler.UpdateConstellationStatusModerator) // Новый маршрут для обновления статуса созвездия модератором
	}

	UserGroup := r.Group("/user")
	{
		UserGroup.GET("/", app.Handler.GetUserByID)
		UserGroup.POST("/register", app.Handler.Register)
		UserGroup.POST("/login", app.Handler.Login)
		UserGroup.POST("/logout", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.Logout)
	}
	addr := fmt.Sprintf("%s:%d", app.Config.ServiceHost, app.Config.ServicePort)
	r.Run(addr)
	log.Println("Server down")
}
