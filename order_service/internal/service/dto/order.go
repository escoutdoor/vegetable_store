package dto

type CreateOrderParams struct {
	UserID      string
	TotalAmount float32
	OrderItems  []CreateOrderItemParams
}

type CreateOrderItemParams struct {
	VegetableID     string
	Weight          float32
	Price           float32
	DiscountedPrice float32

	// Recipient information
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string

	// Address information
	Address string
}

// type CreateRecipientParams struct {
// 	FirstName   string
// 	LastName    string
// 	PhoneNumber string
// 	Email       string
// }
//
// type CreateAddressParams struct {
// 	Address string
// }

type ListOrdersParams struct {
	UserID   string
	Limit    int64
	Offset   int64
	OrderIDs []string
}
