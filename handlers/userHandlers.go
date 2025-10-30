package handlers

import (
	"net/http"

	"github.com/cyrilschreiber3/spotify-playlist-manager/templates/pages"
	"github.com/cyrilschreiber3/spotify-playlist-manager/utils"
	"github.com/gin-gonic/gin"
)

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if exists && userID != "" {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
		component := pages.Login(c, "")
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("session_id", "", -1, "/", "", false, true)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
