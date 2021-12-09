package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yasensim/toyshop/internal/toys"
)

var toysService *ToysService

func Get() *ToysService {
	if toysService == nil {
		toysService = &ToysService{DB: GetToysDataStore()}
		return toysService
	}
	return toysService
}

type ToysService struct {
	DB toys.ToyDatastore
}

func (ts *ToysService) CreateToy(w http.ResponseWriter, r *http.Request) {

	toy := &toys.Toy{}
	json.NewDecoder(r.Body).Decode(toy)

	_, err := ts.DB.FindToy(toy.ProductNumber)

	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ts.DB.CreateToy(toy); err != nil {
		log.Print("error occued CreateToy ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

	var resp = map[string]interface{}{"status": true, "toy": toy}
	json.NewEncoder(w).Encode(resp)
}

func (ts *ToysService) FetchToys(w http.ResponseWriter, r *http.Request) {

	theToys, err := ts.DB.GetAllToys()
	if err != nil {
		log.Print("error occued during FetchToys ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(theToys)
}
