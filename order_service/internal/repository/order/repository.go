package order

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/escoutdoor/vegetable_store/common/pkg/database"
	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/order_service/internal/entity"
	apperrors "github.com/escoutdoor/vegetable_store/order_service/internal/errors"
	"github.com/escoutdoor/vegetable_store/order_service/internal/service/dto"
	"github.com/georgysavva/scany/v2/pgxscan"
)

const (
	// table names
	ordersTableName     = "orders"
	orderItemsTableName = "order_items"
	addressTableName    = "addresses"
	recipientTableName  = "recipients"

	// orders table
	orderIDColumn     = "id"
	orderUserIDColumn = "user_id"
	orderTotalAmount  = "total_amount"

	// order_items table
	orderItemIDColumn              = "id"
	orderItemOrderIDColumn         = "order_id"
	orderItemVegetableIDColumn     = "vegetable_id"
	orderItemWeightColumn          = "weight"
	orderItemPriceColumn           = "price"
	orderItemDiscountedPriceColumn = "discounted_price"
	orderItemAddressIDColumn       = "address_id"
	orderItemRecipientIDColumn     = "recipient_id"

	// addresses table
	addressAddressColumn = "address"

	// recipients table
	recipientFirstNameColumn   = "first_name"
	recipientLastNameColumn    = "last_name"
	recipientPhoneNumberColumn = "phone_number"
	recipientEmailColumn       = "email"

	// default list params values
	defaultLimit  = 10
	defaultOffset = 0
)

type repository struct {
	db database.Client
	qb sq.StatementBuilderType
}

func NewRepository(db database.Client) *repository {
	return &repository{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *repository) CreateOrder(ctx context.Context, in dto.CreateOrderParams) (string, error) {
	sql, args, err := r.qb.Insert(ordersTableName).
		Columns(orderUserIDColumn, orderTotalAmount).
		Values(in.UserID, in.TotalAmount).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return "", buildSQLError(err)
	}

	q := database.Query{
		Name: "order_repository.CreateOrder",
		Sql:  sql,
	}

	var orderID string
	if err := r.db.DB().QueryRowContext(ctx, q, args...).Scan(&orderID); err != nil {
		return "", scanRowError(err)
	}

	for _, oi := range in.OrderItems {
		addressID, err := r.createAddress(ctx, oi.Address)
		if err != nil {
			return "", errwrap.Wrap("create address", err)
		}

		recipientID, err := r.createRecipient(ctx, oi)
		if err != nil {
			return "", errwrap.Wrap("create recipient", err)
		}

		if err := r.createOrderItem(ctx, orderID, addressID, recipientID, oi); err != nil {
			return "", errwrap.Wrap("create order item", err)
		}
	}

	return orderID, nil
}

func (r *repository) createOrderItem(ctx context.Context, orderID, addressID, recipientID string, in dto.CreateOrderItemParams) error {
	sql, args, err := r.qb.Insert(orderItemsTableName).
		Columns(
			orderItemOrderIDColumn,
			orderItemVegetableIDColumn,
			orderItemWeightColumn,
			orderItemPriceColumn,
			orderItemDiscountedPriceColumn,
			orderItemAddressIDColumn,
			orderItemRecipientIDColumn,
		).
		Values(
			orderID,
			in.VegetableID,
			in.Weight,
			in.Price,
			in.DiscountedPrice,
			addressID,
			recipientID,
		).
		ToSql()
	if err != nil {
		return buildSQLError(err)
	}

	q := database.Query{
		Name: "order_repository.createOrderItem",
		Sql:  sql,
	}

	if _, err := r.db.DB().ExecContext(ctx, q, args...); err != nil {
		return executeSQLError(err)
	}

	return nil
}

func (r *repository) createAddress(ctx context.Context, address string) (string, error) {
	sfx := "RETURNING id"
	sql, args, err := r.qb.Insert(addressTableName).
		Columns(addressAddressColumn).
		Values(address).
		Suffix(sfx).
		ToSql()
	if err != nil {
		return "", buildSQLError(err)
	}

	q := database.Query{
		Name: "order_repository.createAddress",
		Sql:  sql,
	}

	var addressID string
	if err := r.db.DB().QueryRowContext(ctx, q, args...).Scan(&addressID); err != nil {
		return "", scanRowError(err)
	}

	return addressID, nil
}

func (r *repository) createRecipient(ctx context.Context, in dto.CreateOrderItemParams) (string, error) {
	sfx := "RETURNING id"
	sql, args, err := r.qb.Insert(recipientTableName).
		Columns(
			recipientFirstNameColumn,
			recipientLastNameColumn,
			recipientPhoneNumberColumn,
			recipientEmailColumn,
		).
		Values(
			in.FirstName,
			in.LastName,
			in.PhoneNumber,
			in.Email,
		).
		Suffix(sfx).
		ToSql()
	if err != nil {
		return "", buildSQLError(err)
	}

	q := database.Query{
		Name: "order_repository.createRecipient",
		Sql:  sql,
	}

	var recipientID string
	if err := r.db.DB().QueryRowContext(ctx, q, args...).Scan(&recipientID); err != nil {
		return "", scanRowError(err)
	}

	return recipientID, nil
}

func (r *repository) GetOrder(ctx context.Context, orderID string) (entity.Order, error) {
	sql, args, err := r.qb.Select(
		"o.id as id",
		"o.user_id",
		"o.total_amount",

		"oi.id as order_item_id",
		"oi.vegetable_id",
		"oi.weight",
		"oi.price",
		"oi.discounted_price",
		"r.first_name",
		"r.last_name",
		"r.phone_number",
		"r.email",
		"a.address",
	).
		From("orders o").
		Join("order_items oi on o.id = oi.order_id").
		Join("addresses a on oi.address_id = a.id").
		Join("recipients r on oi.recipient_id = r.id").
		Where(sq.Eq{"o.id": orderID}).
		ToSql()
	if err != nil {
		return entity.Order{}, buildSQLError(err)
	}

	q := database.Query{
		Name: "order_repository.GetOrder",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return entity.Order{}, executeSQLError(err)
	}
	defer rows.Close()

	var orderRows OrderRows
	if err := pgxscan.ScanAll(&orderRows, rows); err != nil {
		return entity.Order{}, scanRowsError(err)
	}

	if len(orderRows) == 0 {
		return entity.Order{}, apperrors.OrderNotFound(orderID)
	}

	return orderRows.ToServiceEntity(), nil
}

func (r *repository) ListOrders(ctx context.Context, in dto.ListOrdersParams) ([]entity.Order, error) {
	builder := r.qb.Select(
		"o.id as id",
		"o.user_id",
		"o.total_amount",

		"oi.id as order_item_id",
		"oi.vegetable_id",
		"oi.weight",
		"oi.price",
		"oi.discounted_price",
		"r.first_name",
		"r.last_name",
		"r.phone_number",
		"r.email",
		"a.address",
	).
		From("orders o").
		Join("order_items oi on o.id = oi.order_id").
		Join("addresses a on oi.address_id = a.id").
		Join("recipients r on oi.recipient_id = r.id").
		Limit(defaultLimit).
		Offset(defaultOffset)
	if in.Limit > 0 {
		builder = builder.Limit(uint64(in.Limit))
	}
	if in.Offset > 0 {
		builder = builder.Offset(uint64(in.Offset))
	}
	if len(in.OrderIDs) > 0 || in.UserID != "" {
		where := sq.And{}
		if in.UserID != "" {
			where = append(where, sq.Eq{"o.user_id": in.UserID})
		}
		if len(in.OrderIDs) > 0 {
			where = append(where, sq.Eq{"o.id": in.OrderIDs})
		}

		builder = builder.Where(where)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, buildSQLError(err)
	}

	q := database.Query{
		Name: "order_repository.ListOrders",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, executeSQLError(err)
	}
	defer rows.Close()

	var orderRows OrderRows
	if err := pgxscan.ScanAll(&orderRows, rows); err != nil {
		return nil, scanRowsError(err)
	}

	return orderRows.ToServiceEntities(), nil
}
