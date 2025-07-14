package codes

type Code string

const (
	OrderNotFound         Code = "ORDER_NOT_FOUND"
	InsufficientInventory Code = "INSUFFICIENT_INVENTORY"
	VegetablesNotFound    Code = "VEGETABLES_NOT_FOUND"
)
