package main

import (
	"log"
	"net/http"

	"github.com/yasensim/toyshop/toyshop/internal/routes"
)

func main() {
	r := routes.Handlers()

	err := http.ListenAndServe(":5000", r)
	if err != nil {
		log.Fatal(err)
	}
}
