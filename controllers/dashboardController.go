package controllers

import (
	"net/http"

	"github.com/cyrilschreiber3/spotify-playlist-manager/database"
	"github.com/cyrilschreiber3/spotify-playlist-manager/templates/pages"
	"github.com/cyrilschreiber3/spotify-playlist-manager/utils"
	"github.com/gin-gonic/gin"
)

func Dashboard() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id")
		user, err := database.GetUserByID(c.Request.Context(), userID.(string))
		if err != nil {
			utils.RenderTemplate(c, http.StatusInternalServerError, pages.Login(c, "Error fetching user data"))
			return
		}
		utils.RenderTemplate(c, http.StatusOK, pages.Dashboard(c, user))
	}
}
