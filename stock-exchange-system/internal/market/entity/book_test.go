package entity

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestBookTrade(t *testing.T) {
	orderChan := make(chan *Order)
	orderChanOut := make(chan *Order)
	wg := &sync.WaitGroup{}

	book := NewBook(orderChan, orderChanOut, wg)

	// Start the Trade goroutine
	go func() {
		book.Trade()
		close(orderChanOut) // Close the output channel once Trade is done
	}()

	// Add orders
	investor1 := NewInvestor("1")
	investor2 := NewInvestor("2")
	asset := NewAsset("asset1", "Asset 1", 100)

	// Order 1: Sell
	order1 := NewOrder("1", investor1, asset, 10, 5, SellOrderType)
	wg.Add(1)
	orderChan <- order1

	// Order 2: Buy
	order2 := NewOrder("2", investor2, asset, 10, 5, BuyOrderType)
	wg.Add(1)
	orderChan <- order2

	// Close channels after sending orders
	close(orderChan)

	// Wait for the Trade goroutine to finish processing
	wg.Wait()

	// Collect processed orders
	var processedOrders []*Order
	for order := range orderChanOut {
		processedOrders = append(processedOrders, order)
	}

	// Assertions
	assert.Len(t, processedOrders, 2, "There should be 2 processed orders")

	// Verify that the orders have been processed correctly
	assert.Equal(t, Closed, order1.Status, "Order 1 should be closed")
	assert.Equal(t, 0, order1.PendingShares, "Order 1 should have 0 PendingShares")

	assert.Equal(t, Closed, order2.Status, "Order 2 should be closed")
	assert.Equal(t, 0, order2.PendingShares, "Order 2 should have 0 PendingShares")

	// Verify the asset positions
	assert.Equal(t, 5, investor1.GetAssetPosition("asset1").Shares, "Investor 1 should have 5 shares of asset 1")
	assert.Equal(t, 5, investor2.GetAssetPosition("asset1").Shares, "Investor 2 should have 5 shares of asset 1")

	// Optionally verify transactions if necessary
	assert.Len(t, book.Transactions, 1, "There should be 1 transaction")
	transaction := book.Transactions[0]
	assert.Equal(t, 5, transaction.Shares, "Transaction should involve 5 shares")
	assert.Equal(t, 5, transaction.BuyingOrder.Price, "Transaction price should be 5")
	assert.Equal(t, 5, transaction.SellingOrder.Price, "Transaction price should be 5")
}
