package dto

import (
	"invoice-service/internal/domain"
)

type ItemDTO struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

func (i *ItemDTO) ToItemDomain() *domain.Item {
	return domain.NewItem(i.ID, i.Description, i.Price, i.Quantity)
}

func FromItemDomain(i domain.Item) ItemDTO {
	return ItemDTO{
		ID:          i.ID(),
		Description: i.Description(),
		Price:       i.Price(),
		Quantity:    i.Quantity(),
	}
}
