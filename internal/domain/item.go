package domain

type Item struct {
	id          string
	description string
	price       float64
	quantity    int
}

func NewItem(id string, description string, price float64, quantity int) *Item {
	return &Item{
		id:          id,
		description: description,
		price:       price,
		quantity:    quantity,
	}
}

func (i *Item) ID() string {
	return i.id
}

func (i *Item) Description() string {
	return i.description
}

func (i *Item) Price() float64 {
	return i.price
}

func (i *Item) Quantity() int {
	return i.quantity
}
