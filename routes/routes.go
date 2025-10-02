package routes

import (
	"github.com/cyrilschreiber3/spotify-playlist-manager/controllers"
	"github.com/cyrilschreiber3/spotify-playlist-manager/handlers"
	"github.com/cyrilschreiber3/spotify-playlist-manager/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.GET("/", handlers.Index())
	router.Static("/static", "./static")
	router.GET("/login", controllers.Register())
	router.GET("/api/spotifycallback", controllers.SpotifyCallback())

	router.Use(middlewares.AuthMiddleware())

	router.GET("/dashboard", controllers.Dashboard())
	router.GET("/logout", controllers.Logout())
}
