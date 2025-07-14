package vegetable

import (
	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/entity"
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

type Vegetable struct {
	ID              string  `db:"id"`
	Name            string  `db:"name"`
	Weight          float32 `db:"weight"`
	Price           float32 `db:"price"`
	DiscountedPrice float32 `db:"discounted_price"`
}

func (e Vegetable) ToServiceEntity() entity.Vegetable {
	return entity.Vegetable{
		ID:              e.ID,
		Name:            e.Name,
		Weight:          e.Weight,
		Price:           e.Price,
		DiscountedPrice: e.DiscountedPrice,
	}
}

type VegetableList []Vegetable

func (e VegetableList) ToServiceEntity() []entity.Vegetable {
	list := make([]entity.Vegetable, 0, len(e))
	for _, v := range e {
		list = append(list, v.ToServiceEntity())
	}

	return list
}
