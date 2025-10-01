package handlers

import (
	"github.com/gin-gonic/gin"
)

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Redirect user to Spotify authorization URL

	}

}

func SpotifyCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Handle Spotify callback and exchange code for tokens

	}

}
