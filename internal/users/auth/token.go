package auth

import jwt "github.com/golang-jwt/jwt/v4"

//Token struct declaration
type Token struct {
	Name   string
	Email  string
	Region string
	*jwt.StandardClaims
}
