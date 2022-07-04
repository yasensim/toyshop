package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	toysService "github.com/yasensim/toyshop/internal/toys/service"
	"github.com/yasensim/toyshop/internal/users/auth"
	userService "github.com/yasensim/toyshop/internal/users/service"
)

func Handlers() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.Use(CommonMiddleware)
	us := userService.Get()
	ts := toysService.Get()
	av := auth.GetAuthenticator()

	r.HandleFunc("/register", us.CreateUser).Methods("POST")
	r.HandleFunc("/login", us.Login).Methods("POST")

	s := r.PathPrefix("/auth").Subrouter()
	s.Use(av.JwtVerify)
	s.HandleFunc("/toys", ts.GetAllToys).Methods("GET")
	s.HandleFunc("/toys/{id}", ts.GetToy).Methods("GET")
	s.HandleFunc("/toys", ts.CreateToy).Methods("POST")
	s.HandleFunc("/toys/{id}", ts.UpdateToy).Methods("PUT")
	s.HandleFunc("/toys/{id}", ts.DeleteToy).Methods("DELETE")

	return r
}

func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
