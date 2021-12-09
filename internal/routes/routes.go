package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	toysService "github.com/yasensim/toyshop/internal/toys/service"
	"github.com/yasensim/toyshop/internal/users/auth"
	usersService "github.com/yasensim/toyshop/internal/users/service"
)

func Handlers() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.Use(CommonMiddleware)

	us := usersService.Get()
	ts := toysService.Get()

	r.HandleFunc("/register", us.CreateUser).Methods("POST")
	r.HandleFunc("/login", us.Login).Methods("POST")
	jv := auth.GetAuthenticator()
	s := r.PathPrefix("/auth").Subrouter()
	s.Use(jv.JwtVerify)
	s.HandleFunc("/toy", ts.FetchToys).Methods("GET")
	s.HandleFunc("/toy", ts.CreateToy).Methods("POST")
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
