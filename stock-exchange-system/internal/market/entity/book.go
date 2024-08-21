package entity

import (
	"container/heap"
	"sync"
)

const (
	BuyOrderType  = "BUY"
	SellOrderType = "SELL"
)

type Book struct {
	Order         []*Order
	Transactions  []*Transaction
	OrdersChan    chan *Order
	OrdersChanOut chan *Order
	Wg            *sync.WaitGroup
}

func NewBook(orderChan chan *Order, orderChanOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Order:         []*Order{},
		Transactions:  []*Transaction{},
		OrdersChan:    orderChan,
		OrdersChanOut: orderChanOut,
		Wg:            wg,
	}
}

func (b *Book) Trade() {
	buyOrders := make(map[string]*OrderQueue)
	sellOrders := make(map[string]*OrderQueue)
	//buyOrders := NewOrderQueue()
	//sellOrders := NewOrderQueue()

	//heap.Init(buyOrders)
	//heap.Init(sellOrders)

	for order := range b.OrdersChan {
		asset := order.Asset.ID

		if buyOrders[asset] == nil {
			buyOrders[asset] = NewOrderQueue()
			heap.Init(buyOrders[asset])
		}

		if sellOrders[asset] == nil {
			sellOrders[asset] = NewOrderQueue()
			heap.Init(sellOrders[asset])
		}

		switch order.OrderType {
		case BuyOrderType:
			b.processOrder(order, buyOrders[asset], sellOrders[asset], func(o1, o2 *Order) bool {
				return o2.Price <= o1.Price
			})
		case SellOrderType:
			b.processOrder(order, sellOrders[asset], buyOrders[asset], func(o1, o2 *Order) bool {
				return o2.Price >= o1.Price
			})
		}
	}
}

func (b *Book) processOrder(order *Order, primaryQueue, secondaryQueue *OrderQueue, matchCondition func(o1, o2 *Order) bool) {
	primaryQueue.Push(order)
	for primaryQueue.Len() > 0 && secondaryQueue.Len() > 0 && matchCondition(order, secondaryQueue.Orders[0]) {
		matchingOrder := secondaryQueue.Pop().(*Order)
		if matchingOrder.PendingShares > 0 {
			transaction := NewTransaction(order, matchingOrder, order.Shares, matchingOrder.Price)
			b.AddTransaction(transaction, b.Wg)
			matchingOrder.Transactions = append(matchingOrder.Transactions, transaction)
			order.Transactions = append(order.Transactions, transaction)
			b.OrdersChanOut <- matchingOrder
			b.OrdersChanOut <- order
			if matchingOrder.PendingShares > 0 {
				secondaryQueue.Push(matchingOrder)
			}
		}
	}
}

func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	sellingShares := transaction.SellingOrder.PendingShares
	buyingShares := transaction.BuyingOrder.PendingShares

	minShares := sellingShares
	if buyingShares < minShares {
		minShares = buyingShares
	}

	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, -minShares)
	transaction.AddSellOrderPendingShares(-minShares)

	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.BuyingOrder.Asset.ID, minShares)
	transaction.AddBuyOrderPendingShares(-minShares)

	transaction.CalculateTotal(transaction.Shares, transaction.BuyingOrder.Price)

	transaction.CloseBuyOrder()
	transaction.CloseSellOrder()
	b.Transactions = append(b.Transactions, transaction)

}
