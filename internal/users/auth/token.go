package auth

import jwt "github.com/golang-jwt/jwt/v4"

//Token struct declaration
type Token struct {
	ID     string
	Name   string
	Email  string
	Region string
	Org    string
	*jwt.StandardClaims
}
