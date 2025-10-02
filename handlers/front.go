package handlers

import (
	"net/http"

	"github.com/cyrilschreiber3/spotify-playlist-manager/components"
	"github.com/cyrilschreiber3/spotify-playlist-manager/utils"
	"github.com/gin-gonic/gin"
)

func Index() gin.HandlerFunc {
	return func(c *gin.Context) {
		component := components.Index("")
		err := utils.RenderTemplate(c, http.StatusOK, component)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
	}
}
