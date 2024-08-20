package entity

import "sync"

type Book struct {
	Order         []*Order
	Transactions  []*Transaction
	OrderChan     chan *Order
	OrdersChanOut chan *Order
	Wg            *sync.WaitGroup
}

func NewBook(orderChan chan *Order, orderChanOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Order:         []*Order{},
		Transactions:  []*Transaction{},
		OrderChan:     orderChan,
		OrdersChanOut: orderChanOut,
		Wg:            wg,
	}
}
