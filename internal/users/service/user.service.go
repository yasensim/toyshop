package service

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yasensim/toyshop/internal/users"
	"github.com/yasensim/toyshop/internal/users/auth"
)

var usersService *UsersService

func Get() *UsersService {
	if usersService == nil {
		usersService = &UsersService{DB: GetUsersDataStore(), JwtAuth: auth.GetAuthenticator()}
		return usersService
	}
	return usersService
}

type UsersService struct {
	DB      users.UserDatastore
	JwtAuth users.UserAuth
}

func (us *UsersService) Login(w http.ResponseWriter, r *http.Request) {
	user := &users.User{}
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	currUser, err := us.DB.FindUser(user.Email, user.Password)

	if err != nil {
		log.Print("error occued FindUser ", err.Error())
		// delete cookie
		http.SetCookie(w, &http.Cookie{
			Name:       auth.TokenName,
			Value:      "",
			Path:       "/",
			RawExpires: "0",
		})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tokenString, _ := us.JwtAuth.GetTokenForUser(currUser)

	http.SetCookie(w, &http.Cookie{
		Name:       auth.TokenName,
		Value:      tokenString,
		HttpOnly:   true,
		Secure:     false,       // set secure cookie to true for production !!!
		Domain:     "localhost", //change to actual domain
		Path:       "/",
		RawExpires: "0",
	})
	log.Println("User " + currUser.Name + " with email " + currUser.Email + " has loged in!")
	//var resp = map[string]interface{}{"status": true, "access-token": tokenString, "user": currUser}
	var resp = map[string]interface{}{"status": true, "user": currUser.Name}
	json.NewEncoder(w).Encode(resp)
}

func (us *UsersService) CreateUser(w http.ResponseWriter, r *http.Request) {

	user := &users.User{}
	user.ID = uuid.New().String()
	user.Org = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Region = "default"
	user.Active = true
	json.NewDecoder(r.Body).Decode(user)

	/*
		_, err := us.DB.FindUser(user.Email, user.Password)

		if err == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	*/
	if err := us.DB.CreateUser(user); err != nil {
		log.Print("error occued CreateUser ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("User " + user.Name + " with email " + user.Email + " was created!")
	w.WriteHeader(http.StatusCreated)

	var resp = map[string]interface{}{"status": true, "user": user}
	json.NewEncoder(w).Encode(resp)
}

func (us *UsersService) VerifyAuth(w http.ResponseWriter, r *http.Request) {
	exist, token := us.JwtAuth.IsTokenExists(r)
	if !exist {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	usr, err := us.JwtAuth.UserFromToken(token)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usr)
}
