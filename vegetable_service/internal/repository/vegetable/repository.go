package vegetable

import (
	"context"
	"errors"

	"github.com/escoutdoor/vegetable_store/common/pkg/database"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"

	sq "github.com/Masterminds/squirrel"
	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/dto"
	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/entity"
	apperrors "github.com/escoutdoor/vegetable_store/vegetable_service/internal/errors"
	def "github.com/escoutdoor/vegetable_store/vegetable_service/internal/repository"
)

const (
	idColumn              = "id"
	nameColumn            = "name"
	weightColumn          = "weight"
	priceColumn           = "price"
	discountedPriceColumn = "discounted_price"

	tableName = "vegetables"

	defaultLimit  = 50
	defaultOffset = 0
)

type repository struct {
	db database.Client
	qb sq.StatementBuilderType
}

var _ def.VegetableRepository = (*repository)(nil)

func NewRepository(db database.Client) *repository {
	return &repository{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *repository) GetVegetable(ctx context.Context, vegetableID string) (entity.Vegetable, error) {
	sql, args, err := r.qb.Select(idColumn, nameColumn, weightColumn, priceColumn, discountedPriceColumn).
		From(tableName).
		Where(sq.Eq{idColumn: vegetableID}).
		ToSql()
	if err != nil {
		return entity.Vegetable{}, buildSQLError(err)
	}

	q := database.Query{
		Name: "vegetable_repository.GetVegetable",
		Sql:  sql,
	}
	row, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return entity.Vegetable{}, executeSQLError(err)
	}
	defer row.Close()

	var vegetable Vegetable
	if err := pgxscan.ScanOne(&vegetable, row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Vegetable{}, apperrors.VegetableNotFound(vegetableID)
		}

		return entity.Vegetable{}, scanRowError(err)
	}

	return vegetable.ToServiceEntity(), nil
}

func (r *repository) DeleteVegetable(ctx context.Context, vegetableID string) error {
	sql, args, err := r.qb.Delete(tableName).
		Where(sq.Eq{idColumn: vegetableID}).
		ToSql()
	if err != nil {
		return buildSQLError(err)
	}

	q := database.Query{
		Name: "vegetable_repository.DeleteVegetable",
		Sql:  sql,
	}
	if _, err := r.db.DB().ExecContext(ctx, q, args...); err != nil {
		return executeSQLError(err)
	}

	return nil
}

func (r *repository) CreateVegetable(ctx context.Context, in dto.CreateVegetableParams) (string, error) {
	sfx := "RETURNING " + idColumn

	sql, args, err := r.qb.Insert(tableName).
		Columns(nameColumn, weightColumn, priceColumn, discountedPriceColumn).
		Values(in.Name, in.Weight, in.Price, in.DiscountedPrice).
		Suffix(sfx).
		ToSql()
	if err != nil {
		return "", buildSQLError(err)
	}

	q := database.Query{
		Name: "vegetable_repository.CreateVegetable",
		Sql:  sql,
	}
	var vegetableID string
	if err := r.db.DB().QueryRowContext(ctx, q, args...).Scan(&vegetableID); err != nil {
		return "", scanRowError(err)
	}

	return vegetableID, nil
}

func (r *repository) ListVegetables(ctx context.Context, in dto.ListVegetablesParams) ([]entity.Vegetable, error) {
	builder := r.qb.Select(idColumn, nameColumn, weightColumn, priceColumn, discountedPriceColumn).
		From(tableName).
		Limit(defaultLimit).
		Offset(defaultOffset)

	if in.Limit > 0 {
		builder = builder.Limit(uint64(in.Limit))
	}
	if in.Offset > 0 {
		builder = builder.Offset(uint64(in.Offset))
	}
	if len(in.VegetableIDs) > 0 {
		builder = builder.Where(sq.Eq{idColumn: in.VegetableIDs})
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, buildSQLError(err)
	}

	q := database.Query{
		Name: "vegetable_repository.ListVegetables",
		Sql:  sql,
	}
	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, executeSQLError(err)
	}
	defer rows.Close()

	var list VegetableList
	if err := pgxscan.ScanAll(&list, rows); err != nil {
		return nil, scanRowsError(err)
	}

	return list.ToServiceEntity(), nil
}

func (r *repository) UpdateVegetable(ctx context.Context, in dto.VegetableUpdateOperation) error {
	builder := r.qb.Update(tableName).Where(sq.Eq{idColumn: in.ID})

	if in.Name != nil {
		builder = builder.Set(nameColumn, in.Name)
	}
	if in.Weight != nil {
		builder = builder.Set(weightColumn, in.Weight)
	}
	if in.Price != nil {
		builder = builder.Set(priceColumn, in.Price)
	}
	if in.DiscountedPrice != nil {
		builder = builder.Set(discountedPriceColumn, in.DiscountedPrice)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return buildSQLError(err)
	}

	q := database.Query{
		Name: "vegetable_repository.UpdateVegetable",
		Sql:  sql,
	}
	if _, err := r.db.DB().ExecContext(ctx, q, args...); err != nil {
		return executeSQLError(err)
	}

	return nil
}
