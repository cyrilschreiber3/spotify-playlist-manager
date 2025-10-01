package models

type User struct {
	ID           string
	Username     string
	Email        string
	ProfileImage string
}

func NewUser(id, username, email, profileImage string) User {
	return User{
		ID:           id,
		Username:     username,
		Email:        email,
		ProfileImage: profileImage,
	}
}
