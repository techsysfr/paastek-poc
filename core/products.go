package main

var productsDB map[string]*product

type product struct {
	ID       string
	Family   string // `json: "Family, omitempty"`
	Size     string
	Hardware hardware
}

type hardware struct {
	CPU string
	RAM int
}

func getProducts() ([]*product, error) {
	var prds []*product
	for k, v := range productsDB {
		v.ID = k
		prds = append(prds, v)
	}
	return prds, nil
}
func getProduct(id string) (*product, error) {
	return productsDB[id], nil
}

func addProduct(p *product, id string) error {
	productsDB[id] = p
	return nil
}
