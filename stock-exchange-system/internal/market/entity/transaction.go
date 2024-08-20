package entity

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID           string
	SellingOrder *Order
	BuyingOder   *Order
	Shares       int
	Price        float64
	Total        float64
	Datetime     time.Time
}

func NewTransaction(sellingOrder *Order, buyingOrder *Order, shares int, price float64) *Transaction {
	total := float64(shares) * price
	return &Transaction{
		ID:           uuid.New().String(),
		SellingOrder: sellingOrder,
		BuyingOder:   buyingOrder,
		Shares:       shares,
		Price:        price,
		Total:        total,
		Datetime:     time.Now(),
	}
}
