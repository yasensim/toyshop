package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yasensim/toyshop/inventory/internal/inventory"
)

var inventoryService *InventoryService

func Get() *InventoryService {
	if inventoryService == nil {
		inventoryService = &InventoryService{DB: GetInventoryDataStore()}
		return inventoryService
	}
	return inventoryService
}

type InventoryService struct {
	DB inventory.InventoryDatastore
}

func printHeaders(r *http.Request) {
	fmt.Printf("Request at %v\n", time.Now())
	for k, v := range r.Header {
		fmt.Printf("%v: %v\n", k, v)
	}
}
func (is *InventoryService) CreateInventory(w http.ResponseWriter, r *http.Request) {
	printHeaders(r)
	inv := &inventory.Inventory{}
	json.NewDecoder(r.Body).Decode(inv)
	if err := is.DB.CreateInventory(inv); err != nil {
		log.Print("error occured when creating new inventory ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inv)
}

func (is *InventoryService) GetAllInventory(w http.ResponseWriter, r *http.Request) {
	theInvs, err := is.DB.GetAllInventory()
	if err != nil {
		log.Print("error occured when getting all inventory ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(theInvs)
}
