package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/cyrilschreiber3/spotify-playlist-manager/controllers"
	"github.com/cyrilschreiber3/spotify-playlist-manager/database"
	"github.com/cyrilschreiber3/spotify-playlist-manager/routes"
	"github.com/cyrilschreiber3/spotify-playlist-manager/utils"
)

func main() {
	utils.LoadEnv()

	database.Init()
	defer database.Close()

	controllers.InitSpotifyAuth()

	router := gin.Default()
	utils.SetupRouter(router)
	routes.SetupRoutes(router)

	log.Fatal(router.Run(":8080"))
}
