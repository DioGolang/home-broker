package entity

//type Status int
//
//const (
//	Open   Status = iota
//	Closed
//)
//
//func (s Status) String() string{
//	switch s {
//	case Open:
//		return "OPEN"
//	case Closed:
//		return "CLOSED"
//	default:
//		return "Unknown Status"
//	}
//}

type Order struct {
	ID            string
	Investor      *Investor
	Asset         *Asset
	Shares        int
	PendingShares int
	Price         float64
	OrderType     string
	Status        string
	Transactions  []*Transaction
}

func NewOrder(
	orderID string,
	investor *Investor,
	asset *Asset,
	shares int,
	price float64,
	orderType string,
) *Order {
	return &Order{
		ID:           orderID,
		Investor:     investor,
		Asset:        asset,
		Shares:       shares,
		Price:        price,
		OrderType:    orderType,
		Status:       "OPEN",
		Transactions: []*Transactions{},
	}
}
