package obun

import (
	"context"

	dbstore "github.com/otyang/go-dbstore"
	"github.com/uptrace/bun"
)

var _ dbstore.IRepository = (*Repository)(nil)

type (
	SelectCriteria = dbstore.SelectCriteria
	UpdateCriteria = dbstore.UpdateCriteria
	DeleteCriteria = dbstore.DeleteCriteria
)

type Repository struct {
	db bun.IDB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{db: db}
}

func NewRepositoryWithTx(tx bun.Tx) dbstore.IRepository {
	return (&Repository{}).NewWithTx(tx)
}

func (r *Repository) NewWithTx(tx bun.Tx) dbstore.IRepository {
	return &Repository{db: tx}
}

func (r *Repository) Create(ctx context.Context, model any, ignoreDuplicates bool) error {
	if ignoreDuplicates {
		_, err := r.db.NewInsert().Model(model).Ignore().Exec(ctx)
		return err
	}
	_, err := r.db.NewInsert().Model(model).Exec(ctx)
	return err
}

func (r *Repository) CreateBulk(ctx context.Context, modelsPtr any, ignoreDupicates bool) error {
	return r.Create(ctx, modelsPtr, ignoreDupicates)
}

// =========add updateBulk
func (r *Repository) UpdateOneByPK(ctx context.Context, modelPtr any) error {
	_, err := r.db.NewUpdate().Model(modelPtr).WherePK().Exec(ctx)
	return err
}

func (r *Repository) UpdateManyByPK(ctx context.Context, modelPtr any) error {
	_, err := r.db.NewUpdate().Model(modelPtr).WherePK().Bulk().Exec(ctx)
	return err
}

func (r *Repository) UpdateOneWhere(ctx context.Context, modelPtr any, uc ...UpdateCriteria) error {
	q := r.db.NewUpdate().Model(modelPtr)
	for i := range uc {
		if uc[i] == nil {
			continue
		}
		uc[i](q)
	}
	// log.Fatal(q.String())
	_, err := q.Exec(ctx)
	return err
}

func (r *Repository) Upsert(ctx context.Context, modelsPtr any) error {
	_, err := r.db.NewInsert().Model(modelsPtr).On("CONFLICT DO UPDATE").Exec(ctx)
	return err
}

func (r *Repository) DeleteByPK(ctx context.Context, modelPtr any) error {
	_, err := r.db.NewDelete().Model(modelPtr).WherePK().Exec(ctx)
	return err
}

func (r *Repository) DeleteWhere(ctx context.Context, modelPtr any, dc ...DeleteCriteria) error {
	q := r.db.NewDelete().Model(modelPtr)
	for i := range dc {
		if dc[i] == nil {
			continue
		}
		dc[i](q)
	}
	_, err := q.Exec(ctx)
	return err
}

func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context, tx bun.Tx) error) error {
	return r.db.RunInTx(ctx, nil, fn)
}

func (r *Repository) FindOneByPK(ctx context.Context, modelPtr any) error {
	return r.db.NewSelect().Model(modelPtr).WherePK().Limit(1).Scan(ctx)
}

func (r *Repository) FindOneWhere(ctx context.Context, modelPtr any, sc ...SelectCriteria) error {
	q := r.db.NewSelect().Model(modelPtr)

	for i := range sc {
		if sc[i] == nil {
			continue
		}
		sc[i](q)
	}

	return q.Limit(1).Scan(ctx)
}

func (r *Repository) FindManyWhere(ctx context.Context, modelPtr any, opt dbstore.PaginationOption, sc ...SelectCriteria) error {
	q := r.db.NewSelect().Model(modelPtr)
	for i := range sc {
		if sc[i] == nil {
			continue
		}
		sc[i](q)
	}

	if opt == nil {
		if err := q.Scan(ctx); err != nil {
			return err
		}
		return nil
	}

	o := dbstore.PaginationParams{}
	if err := opt(&o); err != nil {
		return err
	}

	// decide if sort defined or predefined
	if o.DirectionNextPage {
		q = q.OrderExpr(o.CursorColumn + " ASC").Limit(o.Limit)
		if o.CursorValue != "" {
			q = q.Where("? >= ?", bun.Ident(o.CursorColumn), o.CursorValue)
		}
		return q.Scan(ctx)
	}

	q = q.OrderExpr(o.CursorColumn + " DESC").Limit(o.Limit)
	if o.CursorValue != "" {
		q = q.Where("? <= ?", bun.Ident(o.CursorColumn), o.CursorValue)
	}
	return q.Scan(ctx)
}
