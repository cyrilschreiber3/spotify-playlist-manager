package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/cyrilschreiber3/spotify-playlist-manager/components"
	"github.com/cyrilschreiber3/spotify-playlist-manager/controllers"
	"github.com/cyrilschreiber3/spotify-playlist-manager/database"
	"github.com/cyrilschreiber3/spotify-playlist-manager/utils"
	"github.com/gin-gonic/gin"
)

func AuthenticateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil || sessionID == "" {
			session, err := controllers.NewSession(c)
			if err != nil {
				log.Println("Error creating new session:", err)
				utils.RenderTemplate(c, http.StatusInternalServerError, components.Login("Error creating session"))
				return
			}
			c.Set("session_id", session.SessionID)
			c.SetCookie("session_id", session.SessionID, 3600*24, "/", "", false, true)
			c.Next()

			return
		}

		session, err := database.GetSessionByID(c.Request.Context(), sessionID)
		if err != nil {
			log.Println("Invalid session ID:", err)
			c.Redirect(302, "/login")
			c.Abort()
			return
		}
		c.Set("session_id", session.SessionID)

		if session.ExpiresAt.Before(time.Now()) {
			log.Println("Session expired for session ID:", sessionID)
			c.SetCookie("session_id", "", -1, "/", "", false, true)
			c.Redirect(302, "/login")
			c.Abort()
			return
		}

		if session.UserID == "" {
			c.Next()
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

func AllowAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.Get("username")
		if exists && username != "" {
			c.Next()
			return
		}

		c.Redirect(302, "/login")
		c.Abort()
	}
}
