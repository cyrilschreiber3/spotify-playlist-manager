package models

import "time"

type User struct {
	ID           string `bson:"_id,omitempty" json:"_id,omitempty"`
	UserID       string `bson:"user_id" json:"user_id"`
	Username     string `bson:"username" json:"username" validate:"required,min=2,max=100"`
	Email        string `bson:"email" json:"email" validate:"required,email"`
	ProfileImage string `bson:"profile_image" json:"profile_image"`
}

type UserSession struct {
	UserID       string    `bson:"user_id" json:"user_id" validate:"required"`
	SessionID    string    `bson:"session_id" json:"session_id" validate:"required"`
	SpotifyToken string    `bson:"spotify_token" json:"spotify_token"`
	ExpiresAt    time.Time `bson:"expires_at" json:"expires_at" validate:"required"`
}

func NewUser(id, username, email, profileImage string) User {
	return User{
		ID:           id,
		Username:     username,
		Email:        email,
		ProfileImage: profileImage,
	}
}
