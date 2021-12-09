package service

import (
	"context"
	"database/sql"
	"log"
	"time"

	database "github.com/yasensim/toyshop/inventory/internal/db"
	"github.com/yasensim/toyshop/inventory/internal/inventory"
)

type InventoryDB struct {
	*sql.DB
}

func GetInventoryDataStore() inventory.InventoryDatastore {
	return &InventoryDB{database.Get()}
}

func (db *InventoryDB) CreateInventory(inv *inventory.Inventory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := db.ExecContext(ctx, "insert into inventory (product_number, quantity) values (?, ?)",
		inv.ProductNumber, inv.Quantity)
	if err != nil {
		return err
	}
	id, e := result.LastInsertId()
	if e != nil {
		return e
	}

	inv.ID = uint(id)

	return nil
}

func (db *InventoryDB) GetAllInventory() ([]inventory.Inventory, error) {
	var theInv []inventory.Inventory

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, "select id, product_number, quantity from inventory")
	if err != nil {
		log.Print("x error occured when getting all inventory ", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var inv inventory.Inventory
		rows.Scan(&inv.ID, &inv.ProductNumber, &inv.Quantity)
		theInv = append(theInv, inv)
	}
	return theInv, nil
}
