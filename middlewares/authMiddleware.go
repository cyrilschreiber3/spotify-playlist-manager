package middlewares

import (
	"log"

	"github.com/cyrilschreiber3/spotify-playlist-manager/database"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil || sessionID == "" {
			log.Println("No ession coockie found:", err)
			c.Redirect(302, "/login")
			c.Abort()
			return
		}

		session, err := database.GetSessionByID(c.Request.Context(), sessionID)
		if err != nil {
			log.Println("Invalid session ID:", err)
			c.Redirect(302, "/login")
			c.Abort()
			return
		}

		user, err := database.GetUserByID(c.Request.Context(), session.UserID)
		if err != nil {
			log.Println("User not found for session:", err)
			c.SetCookie("session_id", "", -1, "/", "", false, true)
			c.Redirect(302, "/login")
			c.Abort()
			return
		}

		c.Set("user_id", session.UserID)
		c.Set("username", user.Username)
		c.Set("email", user.Email)
		c.Set("profile_image", user.ProfileImage)
		c.Set("spotify_token", user.SpotifyToken)

		c.Next()
	}
}
