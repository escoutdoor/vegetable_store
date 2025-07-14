package order

import (
	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/order_service/internal/entity"
)

func buildSQLError(err error) error {
	return errwrap.Wrap("build sql", err)
}

func executeSQLError(err error) error {
	return errwrap.Wrap("execute sql", err)
}

func scanRowError(err error) error {
	return errwrap.Wrap("scan row", err)
}

func scanRowsError(err error) error {
	return errwrap.Wrap("scan rows", err)
}

type OrderRow struct {
	ID          string  `db:"id"`
	UserID      string  `db:"user_id"`
	TotalAmount float32 `db:"total_amount"`

	OrderItemID     string  `db:"order_item_id"`
	VegetableID     string  `db:"vegetable_id"`
	Weight          float32 `db:"weight"`
	Price           float32 `db:"price"`
	DiscountedPrice float32 `db:"discounted_price"`

	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	PhoneNumber string `db:"phone_number"`
	Email       string `db:"email"`

	Address string `db:"address"`
}

type OrderRows []OrderRow

func (rows OrderRows) ToServiceEntities() []entity.Order {
	m := make(map[string]OrderRows)
	for _, o := range rows {
		m[o.ID] = append(m[o.ID], o)
	}

	orders := make([]entity.Order, 0, len(m))
	for id, items := range m {
		orderItems := []entity.OrderItem{}
		for _, i := range items {
			orderItems = append(orderItems, entity.OrderItem{
				ID:              i.ID,
				VegetableID:     i.VegetableID,
				Weight:          i.Weight,
				Price:           i.Price,
				DiscountedPrice: i.DiscountedPrice,
				AddressInfo: entity.AddressInfo{
					Address: i.Address,
				},
				RecipientInfo: entity.RecipientInfo{
					FirstName:   i.FirstName,
					LastName:    i.LastName,
					PhoneNumber: i.PhoneNumber,
					Email:       i.Email,
				},
			})
		}

		order := entity.Order{
			ID:          id,
			UserID:      items[0].UserID,
			TotalAmount: items[0].TotalAmount,
			OrderItems:  orderItems,
		}
		orders = append(orders, order)
	}

	return orders
}

func (r OrderRows) ToServiceEntity() entity.Order {
	list := r.ToServiceEntities()

	return list[0]
}
