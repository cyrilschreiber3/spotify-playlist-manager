package routes

import (
	"github.com/cyrilschreiber3/spotify-playlist-manager/controllers"
	"github.com/cyrilschreiber3/spotify-playlist-manager/handlers"
	"github.com/cyrilschreiber3/spotify-playlist-manager/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.Static("/static", "./static")

	router.Use(middlewares.AuthenticateUser())

	router.GET("/login", handlers.Login())
	router.GET("/logout", handlers.Logout())
	router.GET("/api/login", controllers.Register())
	router.GET("/api/spotifycallback", controllers.SpotifyCallback())

	router.GET("/", middlewares.AllowAuthenticated(), controllers.Dashboard())
}
