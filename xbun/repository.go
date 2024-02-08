package xbun

import (
	"context"

	dbstore "github.com/otyang/go-dbstore"
	"github.com/uptrace/bun"
)

type (
	SelectCriteria func(*bun.SelectQuery) *bun.SelectQuery
	UpdateCriteria func(*bun.UpdateQuery) *bun.UpdateQuery
	DeleteCriteria func(*bun.DeleteQuery) *bun.DeleteQuery
)

// Creates a single record. Returns error if it already exists.
// To silently discard the error set ignoreDuplicate errors to true.
// Note ignoring duplicates doesnt mean the data will be inserted. it
// just ensures the query exits silently
func Create[T any](ctx context.Context, db bun.IDB, model *T, ignoreDuplicates bool) error {
	if ignoreDuplicates {
		_, err := db.NewInsert().Model(model).Ignore().Exec(ctx)
		return err
	}
	_, err := db.NewInsert().Model(model).Exec(ctx)
	return err
}

// Creates a multiple record. ignore duplocate runs SQL on conflict ignore duplicate
func CreateBulk[T any](ctx context.Context, db bun.IDB, model *[]T, ignoreDuplicates bool) error {
	if ignoreDuplicates {
		_, err := db.NewInsert().Model(model).Ignore().Exec(ctx)
		return err
	}
	_, err := db.NewInsert().Model(model).Exec(ctx)
	return err
}

func FindOneByPK[T any](ctx context.Context, db bun.IDB, modelPtr *T) error {
	return db.NewSelect().Model(modelPtr).WherePK().Limit(1).Scan(ctx)
}

func FindOneWhere[T any](ctx context.Context, db bun.IDB, modelPtr *T, sc ...SelectCriteria) error {
	q := db.NewSelect().Model(modelPtr)

	for i := range sc {
		if sc[i] == nil {
			continue
		}
		sc[i](q)
	}

	return q.Limit(1).Scan(ctx)
}

func FindManyWhere[T any](ctx context.Context, db bun.IDB, opt dbstore.PaginationOption, sc ...SelectCriteria) ([]T, error) {
	var modelsPtr []T

	q := db.NewSelect().Model(&modelsPtr)
	for i := range sc {
		if sc[i] == nil {
			continue
		}
		sc[i](q)
	}

	if opt == nil {
		if err := q.Scan(ctx); err != nil {
			return nil, err
		}
		return modelsPtr, nil
	}

	o := dbstore.PaginationParams{}
	if err := opt(&o); err != nil {
		return nil, err
	}

	if o.DirectionNextPage {
		q = q.OrderExpr(o.CursorColumn + " ASC").Limit(o.Limit)
		if o.CursorValue != "" {
			q = q.Where("? >= ?", bun.Ident(o.CursorColumn), o.CursorValue)
		}
		err := q.Scan(ctx)

		return modelsPtr, err
	}

	q = q.OrderExpr(o.CursorColumn + " DESC").Limit(o.Limit)
	if o.CursorValue != "" {
		q = q.Where("? <= ?", bun.Ident(o.CursorColumn), o.CursorValue)
	}
	err := q.Scan(ctx)

	return modelsPtr, err
}

func UpdateOneByPK[T any](ctx context.Context, db bun.IDB, modelPtr *T) error {
	_, err := db.NewUpdate().Model(modelPtr).WherePK().Exec(ctx)
	return err
}

func UpdateManyByPK[T any](ctx context.Context, db bun.IDB, modelPtr *[]T) error {
	_, err := db.NewUpdate().Model(modelPtr).WherePK().Bulk().Exec(ctx)
	return err
}

func UpdateOneWhere[T any](ctx context.Context, db bun.IDB, modelPtr *T, uc ...UpdateCriteria) error {
	q := db.NewUpdate().Model(modelPtr)
	for i := range uc {
		if uc[i] == nil {
			continue
		}
		uc[i](q)
	}
	_, err := q.Exec(ctx)
	return err
}

func Upsert[T any](ctx context.Context, db bun.IDB, modelsPtr *T) error {
	_, err := db.NewInsert().Model(modelsPtr).On("CONFLICT DO UPDATE").Exec(ctx)
	return err
}

func DeleteByPK[T any](ctx context.Context, db bun.IDB, modelPtr *T) error {
	_, err := db.NewDelete().Model(modelPtr).WherePK().Exec(ctx)
	return err
}

func DeleteWhere[T any](ctx context.Context, db bun.IDB, modelPtr *T, dc ...DeleteCriteria) error {
	q := db.NewDelete().Model(modelPtr)
	for i := range dc {
		if dc[i] == nil {
			continue
		}
		dc[i](q)
	}
	_, err := q.Exec(ctx)
	return err
}

func Transaction(ctx context.Context, db *bun.DB, fn func(ctx context.Context, tx bun.Tx) error) error {
	return db.RunInTx(ctx, nil, fn)
}
