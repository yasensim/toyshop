package users

import (
	"net/http"
	"time"
)

// This should actually be pulled into a separate package
// since used from multiple locations
type User struct {
	ID        string    `json:"_id" dynamodbav:"id"`
	Name      string    `json:"name" dynamodbav:"name"`
	Email     string    `json:"email" dynamodbav:"email"`
	Password  string    `json:"password" dynamodb:"password"`
	Org       string    `json:"org" dynamodbav:"org"`
	Region    string    `json:"region" dynamodbav:"reg"`
	Active    bool      `json:"active" dynamodbav:"act"`
	Info      string    `json:"info" dynamodbav:"info"`
	CreatedAt time.Time `json:"createdAt" dynamodbav:"cr"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"upd"`
}

type UserDatastore interface {
	CreateUser(user *User) error
	FindUser(email, password string) (*User, error)
}

type UserAuth interface {
	IsTokenExists(r *http.Request) (bool, string)
	IsUserTokenValid(token string) bool
	UserFromToken(tokenString string) (*User, error)
	GetTokenForUser(user *User) (string, error)
}
