package controllers

import (
	"context"
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

func NewSession(c context.Context) (models.UserSession, error) {
	var session models.UserSession

	session.SessionID = uuid.NewString()
	session.UserID = ""
	session.ExpiresAt = time.Now().Add(24 * time.Hour)

	_, err := database.CreateSession(c, &session)
	if err != nil {
		return models.UserSession{}, err
	}

	return session, err
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		component := components.Login(c, "")
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.MustGet("session_id").(string)

		registerURL := auth.AuthURL(sessionID)

		c.Writer.Header().Set("HX-Redirect", registerURL)

	}

}

func SpotifyCallback() gin.HandlerFunc {
	return func(c *gin.Context) {

		sessionID := c.MustGet("session_id").(string)
		if c.Query("state") != sessionID {
			utils.RenderTemplate(c, http.StatusForbidden, components.Login(c, "Invalid session state"))
			return
		}

		token, err := auth.Token(c.Request.Context(), c.Query("state"), c.Request)
		if err != nil {
			utils.RenderTemplate(c, http.StatusForbidden, components.Login(c, "Error getting token"))
			return
		}

		client := spotify.New(auth.Client(c.Request.Context(), token))

		spotifyUser, err := client.CurrentUser(c.Request.Context())
		if err != nil {
			utils.RenderTemplate(c, http.StatusForbidden, components.Login(c, "Error fetching user data from Spotify"))
			return
		}

		var user models.User
		user.UserID = spotifyUser.ID
		user.Username = spotifyUser.DisplayName
		user.Email = spotifyUser.Email
		if len(spotifyUser.Images) > 0 {
			user.ProfileImage = spotifyUser.Images[0].URL
		} else {
			user.ProfileImage = fmt.Sprintf("https://ui-avatars.com/api/?name=%c", user.Username[0])
		}
		user.SpotifyToken = token.AccessToken

		_, err = database.GetUserByID(c.Request.Context(), user.UserID)
		if err != nil && err != mongo.ErrNoDocuments {
			log.Println("Database error:", err)
			utils.RenderTemplate(c, http.StatusInternalServerError, components.Login(c, "Database error"))
			return
		}

		if err == mongo.ErrNoDocuments {
			_, err = database.CreateUser(c.Request.Context(), &user)
			if err != nil {
				utils.RenderTemplate(c, http.StatusInternalServerError, components.Login(c, "Error creating user"))
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
				utils.RenderTemplate(c, http.StatusInternalServerError, components.Login(c, "Error updating user"))
				return
			}
		}

		session, err := database.GetSessionByID(c.Request.Context(), sessionID)
		if err != nil {
			utils.RenderTemplate(c, http.StatusInternalServerError, components.Login(c, "Error fetching session"))
			return
		}

		session.UserID = user.UserID
		session.ExpiresAt = time.Now().Add(24 * time.Hour)

		_, err = database.UpdateSessionByID(c.Request.Context(), session.SessionID, bson.M{
			"user_id":    session.UserID,
			"expires_at": session.ExpiresAt,
		})
		if err != nil {
			utils.RenderTemplate(c, http.StatusInternalServerError, components.Login(c, "Error updating session"))
			return
		}

		c.SetCookie("session_id", session.SessionID, 3600*24, "/", "", false, true)

		c.Redirect(http.StatusTemporaryRedirect, "/")

	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("session_id", "", -1, "/", "", false, true)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
