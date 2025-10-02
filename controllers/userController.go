package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cyrilschreiber3/spotify-playlist-manager/components"
	"github.com/cyrilschreiber3/spotify-playlist-manager/database"
	"github.com/cyrilschreiber3/spotify-playlist-manager/models"
	"github.com/cyrilschreiber3/spotify-playlist-manager/utils"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/google/uuid"
)

var redirectURL string
var auth *spotifyauth.Authenticator

func InitSpotifyAuth() {
	scopes := []string{
		spotifyauth.ScopeUserReadEmail,
		spotifyauth.ScopePlaylistReadPrivate,
		spotifyauth.ScopePlaylistReadCollaborative,
		spotifyauth.ScopePlaylistModifyPublic,
		spotifyauth.ScopePlaylistModifyPrivate,
	}
	redirectURL = fmt.Sprintf("%sapi/spotifycallback", utils.GetEnv("APP_URL", "http://localhost:8080/"))
	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURL), spotifyauth.WithScopes(scopes...))
}

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var session models.UserSession

		session.SessionID = uuid.NewString()
		session.UserID = ""
		session.ExpiresAt = time.Now()

		_, err := database.CreateSession(c.Request.Context(), &session)
		if err != nil {
			err = utils.RenderTemplate(c, http.StatusInternalServerError, components.Index("Error creating session"))
			if err != nil {
				c.Status(http.StatusInternalServerError)
			}
			return
		}

		registerURL := auth.AuthURL(session.SessionID)

		c.Redirect(http.StatusTemporaryRedirect, registerURL)

	}

}

func SpotifyCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := auth.Token(c.Request.Context(), c.Query("state"), c.Request)
		if err != nil {
			err = utils.RenderTemplate(c, http.StatusForbidden, components.Index("Error getting token"))
			if err != nil {
				c.Status(http.StatusInternalServerError)
			}
			return
		}

		client := spotify.New(auth.Client(c.Request.Context(), token))

		spotifyUser, err := client.CurrentUser(c.Request.Context())
		if err != nil {
			err = utils.RenderTemplate(c, http.StatusForbidden, components.Index("Error fetching user data from Spotify"))
			if err != nil {
				c.Status(http.StatusInternalServerError)
			}
			return
		}

		var user models.User
		user.UserID = spotifyUser.ID
		user.Username = spotifyUser.DisplayName
		user.Email = spotifyUser.Email
		if len(spotifyUser.Images) > 0 {
			user.ProfileImage = spotifyUser.Images[0].URL
		}
		user.SpotifyToken = token.AccessToken

		_, err = database.GetUserByID(c.Request.Context(), user.UserID)
		if err != nil && err != mongo.ErrNoDocuments {
			err = utils.RenderTemplate(c, http.StatusInternalServerError, components.Index("Database error"))
			log.Println("Database error:", err)
			if err != nil {
				c.Status(http.StatusInternalServerError)
			}
			return
		}

		if err == mongo.ErrNoDocuments {
			_, err = database.CreateUser(c.Request.Context(), &user)
			if err != nil {
				err = utils.RenderTemplate(c, http.StatusInternalServerError, components.Index("Error creating user"))
				if err != nil {
					c.Status(http.StatusInternalServerError)
				}
				return
			}
		} else {
			updatedUser := bson.M{
				"username":      user.Username,
				"email":         user.Email,
				"profile_image": user.ProfileImage,
				"spotify_token": user.SpotifyToken,
			}

			_, err = database.UpdateUserByID(c.Request.Context(), user.UserID, updatedUser)
			if err != nil {
				err = utils.RenderTemplate(c, http.StatusInternalServerError, components.Index("Error updating user"))
				if err != nil {
					c.Status(http.StatusInternalServerError)
				}
				return
			}
		}

		sessionID := c.Query("state")
		session, err := database.GetSessionByID(c.Request.Context(), sessionID)
		if err != nil {
			err = utils.RenderTemplate(c, http.StatusInternalServerError, components.Index("Error fetching session"))
			if err != nil {
				c.Status(http.StatusInternalServerError)
			}
			return
		}

		session.UserID = user.UserID
		session.ExpiresAt = time.Now().Add(24 * time.Hour)

		_, err = database.UpdateSessionByID(c.Request.Context(), session.SessionID, bson.M{
			"user_id":    session.UserID,
			"expires_at": session.ExpiresAt,
		})
		if err != nil {
			err = utils.RenderTemplate(c, http.StatusInternalServerError, components.Index("Error updating session"))
			if err != nil {
				c.Status(http.StatusInternalServerError)
			}
			return
		}

		c.SetCookie("session_id", session.SessionID, 3600*24, "/", "", false, true)

		c.Redirect(http.StatusTemporaryRedirect, "/dashboard")

	}
}

func Dashboard() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			err := utils.RenderTemplate(c, http.StatusInternalServerError, components.Index("User ID not found in context"))
			if err != nil {
				c.Status(http.StatusInternalServerError)
			}
			return
		}
		user, err := database.GetUserByID(c.Request.Context(), userID.(string))
		if err != nil {
			err = utils.RenderTemplate(c, http.StatusInternalServerError, components.Index("Error fetching user data"))
			if err != nil {
				c.Status(http.StatusInternalServerError)
			}
			return
		}
		err = utils.RenderTemplate(c, http.StatusOK, components.Dashboard(user))
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("session_id", "", -1, "/", "", false, true)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
