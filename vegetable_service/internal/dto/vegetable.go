package dto

type ListVegetablesParams struct {
	Limit        int
	Offset       int
	VegetableIDs []string
}

type CreateVegetableParams struct {
	Name            string
	Weight          float32
	Price           float32
	DiscountedPrice float32
}

type VegetableUpdateOperation struct {
	ID     string
	Name   *string
	Weight *float32

	Price           *float32
	DiscountedPrice *float32
}

type VegetablesBatchUpdateOperation struct {
	Requests []VegetableUpdateOperation
}
