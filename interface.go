package dbstore

import (
	"context"

	"github.com/uptrace/bun"
)

type (
	SelectCriteria func(*bun.SelectQuery) *bun.SelectQuery
	UpdateCriteria func(*bun.UpdateQuery) *bun.UpdateQuery
	DeleteCriteria func(*bun.DeleteQuery) *bun.DeleteQuery
)

// IRepository is an interface for generic implementation of repository patterns.
// It allows for writing less code and easily perform CRUD operations.
type IRepository interface {
	// Create inserts a single record into the database.
	// It optionally suppresses duplicate key errors.
	Create(ctx context.Context, modelPtr any, suppressDuplicateError bool) error

	// CreateBulk inserts multiple records into the database in a single transaction.
	// It optionally suppresses duplicate key errors.
	CreateBulk(ctx context.Context, modelsPtr any, suppressDuplicateError bool) error

	// FindOneByPK retrieves a single record by its primary key.
	FindOneByPK(ctx context.Context, modelPtr any) error

	// FindOneWhere retrieves a single record matching the specified criteria.
	FindOneWhere(ctx context.Context, modelPtr any, sc ...SelectCriteria) error

	// FindManyWhere retrieves multiple records matching the specified criteria.
	// It supports pagination using the provided PaginationOption.
	FindManyWhere(ctx context.Context, modelPtr any, opt PaginationOption, sc ...SelectCriteria) error

	// UpdateOneByPK updates a single record by its primary key.
	UpdateOneByPK(ctx context.Context, modelsPtr any) error

	// UpdateManyByPK updates multiple records by their primary keys.
	UpdateManyByPK(ctx context.Context, modelsPtr any) error

	// UpdateOneWhere updates a single record matching the specified criteria.
	UpdateOneWhere(ctx context.Context, modelPtr any, uc ...UpdateCriteria) error

	// Upsert inserts a record if it doesn't exist, or updates it if it does.
	Upsert(ctx context.Context, modelsPtr any) error

	// DeleteByPK deletes a single record by its primary key.
	DeleteByPK(ctx context.Context, modelsPtr any) error

	// DeleteWhere deletes multiple records matching the specified criteria.
	DeleteWhere(ctx context.Context, modelPtr any, dc ...DeleteCriteria) error

	// NewWithTx creates a new repository instance using an existing bun.Tx transaction.
	NewWithTx(tx bun.Tx) IRepository

	// Transaction executes a function within a database transaction.
	// It Simplifies transactions handling, by automatically:
	//
	// 	- starting a transaction
	// 	- rolling back the transaction if an error occurs
	// 	- And finally commiting the transaction if no error.
	Transaction(ctx context.Context, fn func(ctx context.Context, tx bun.Tx) error) error
}
