package entity

type Order struct {
	ID          string
	UserID      string
	TotalAmount float32
	OrderItems  []OrderItem
}

type OrderItem struct {
	ID              string
	VegetableID     string
	Weight          float32
	Price           float32
	DiscountedPrice float32

	AddressInfo   AddressInfo
	RecipientInfo RecipientInfo
}

type RecipientInfo struct {
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
}

type AddressInfo struct {
	Address string
}
