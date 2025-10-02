package models

import "time"

type User struct {
	UserID       string `bson:"user_id" json:"user_id" validate:"required"`
	Username     string `bson:"username" json:"username" validate:"required,min=2,max=100"`
	Email        string `bson:"email" json:"email" validate:"email"`
	ProfileImage string `bson:"profile_image" json:"profile_image"`
	SpotifyToken string `bson:"spotify_token" json:"spotify_token"`
}

type UserSession struct {
	UserID    string    `bson:"user_id" json:"user_id"`
	SessionID string    `bson:"session_id" json:"session_id" validate:"required"`
	ExpiresAt time.Time `bson:"expires_at" json:"expires_at" validate:"required"`
}
