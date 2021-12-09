package toys

type Toy struct {
	ID            uint   `json:"id"`
	ProductNumber string `json:"productnumber"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Cost          int    `json:"cost"`
}

type ToyDatastore interface {
	CreateToy(toy *Toy) error
	FindToy(productnumber string) (*Toy, error)
	GetAllToys() ([]Toy, error)
}
