package auth

import jwt "github.com/golang-jwt/jwt/v4"

//Token struct declaration
type Token struct {
	UserID uint
	Name   string
	Email  string
	*jwt.StandardClaims
}
