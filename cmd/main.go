package main

import (
	"log"
	"net/http"

	"github.com/yasensim/toyshop/internal/routes"
)

func main() {
	r := routes.Handlers()

	err := http.ListenAndServe(":5001", r)
	if err != nil {
		log.Fatal(err)
	}
}
