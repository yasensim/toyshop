package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	database "github.com/yasensim/toyshop/internal/db"
	"github.com/yasensim/toyshop/internal/toys"
)

type ToyDB struct {
	*sql.DB
}

func GetToysDataStore() toys.ToyDatastore {
	return &ToyDB{database.Get()}
}

func (db *ToyDB) CreateToy(toy *toys.Toy) error {
	if toy.Name == "" || toy.ProductNumber == "" {
		return errors.New("toy service repo - cannot have empty fields")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	result, err := db.ExecContext(ctx, "insert into toys (id, productnumber, name, description, cost) values (?, ?, ?, ?, ?)",
		toy.ID, toy.ProductNumber, toy.Name, toy.Description, toy.Cost)

	if err != nil {
		return err
	}

	id, e := result.LastInsertId()
	if e != nil {
		return e
	}

	toy.ID = uint(id)

	return nil
}

func (db *ToyDB) FindToy(productnumber string) (*toys.Toy, error) {
	toy := &toys.Toy{}

	if productnumber == "" {
		return nil, errors.New("toy service repo - cannot have empty product number")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	row := db.QueryRowContext(ctx, "select id, productnumber, name, description, cost from toys where productnumber = ?", productnumber)
	err := row.Scan(&toy.ID, &toy.ProductNumber, &toy.Name, &toy.Description, toy.Cost)

	if err == sql.ErrNoRows {
		return nil, err
	}

	return toy, nil
}

func (db *ToyDB) GetAllToys() ([]toys.Toy, error) {
	var theToys []toys.Toy

	rows, err := db.Query("select id, productnumber, name, description, cost from toys")

	if err != nil {
		log.Print("error occued during toys fetch ", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var toy toys.Toy
		rows.Scan(&toy.ID, &toy.ProductNumber, &toy.Name, &toy.Description, toy.Cost)
		theToys = append(theToys, toy)
	}
	return theToys, nil
}
