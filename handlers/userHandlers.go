package handlers

import (
	"net/http"

	"github.com/cyrilschreiber3/spotify-playlist-manager/components"
	"github.com/cyrilschreiber3/spotify-playlist-manager/utils"
	"github.com/gin-gonic/gin"
)

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		component := components.Login("")
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("session_id", "", -1, "/", "", false, true)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
