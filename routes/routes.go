package routes

import (
	"github.com/cyrilschreiber3/spotify-playlist-manager/handlers"
	"github.com/cyrilschreiber3/spotify-playlist-manager/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.GET("/", handlers.Index())
	router.Static("/static", "./static")
	router.GET("/login", handlers.Login())

	router.Use(middlewares.AuthMiddleware())

	router.GET("/dashboard", handlers.Dashboard())
	router.GET("/api/spotifycallback", handlers.SpotifyCallback())
}
