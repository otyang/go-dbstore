package seeder

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

var (
	ErrCreateTablesPrefix     = "create table error: %w"
	ErrDropTablesPrefix       = "drop table error: %w"
	ErrDropCreateTablesPrefix = "drop and create tables error: %w"
)

type Seeder struct {
	db *bun.DB
}

func NewSeeder(db *bun.DB) *Seeder {
	return &Seeder{db: db}
}

func (sm *Seeder) RegisterModels(ctx context.Context, models []any, intermediaryModels []any) {
	sm.db.RegisterModel(intermediaryModels...)
	sm.db.RegisterModel(models...)
}

func (sm *Seeder) CreateTables(ctx context.Context, models []any, intermediaryModels []any) error {
	err := sm.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		ms := append(models, intermediaryModels...)
		for _, model := range ms {
			if _, err := tx.NewCreateTable().Model(model).Exec(ctx); err != nil {
				return fmt.Errorf(ErrCreateTablesPrefix, err)
			}
		}

		return nil
	})
	return err
}

func (sm *Seeder) DropTables(ctx context.Context, models []any, intermediaryModels []any) error {
	err := sm.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		ms := append(models, intermediaryModels...)
		for _, model := range ms {
			if _, err := tx.NewDropTable().Model(model).Cascade().IfExists().Exec(ctx); err != nil {
				return fmt.Errorf(ErrDropTablesPrefix, err)
			}
		}
		return nil
	})

	return err
}

func (sm *Seeder) DropAndCreateTables(ctx context.Context, models []any, intermediaryModels []any) error {
	err := sm.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		ms := append(models, intermediaryModels...)

		for _, model := range ms {
			if _, err := tx.NewDropTable().Model(model).Cascade().IfExists().Exec(ctx); err != nil {
				return fmt.Errorf(ErrDropCreateTablesPrefix, err)
			}
		}
		for _, model := range ms {
			if _, err := tx.NewCreateTable().Model(model).Exec(ctx); err != nil {
				return fmt.Errorf(ErrDropCreateTablesPrefix, err)
			}
		}
		return nil
	})

	return err
}

func (sm *Seeder) CreateIndex(ctx context.Context, modelPtr any, indexName string, indexColumn string) error {
	err := sm.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := sm.db.NewCreateIndex().Model(modelPtr).Index(indexName).Column(indexColumn).Exec(ctx)
		return err
	})
	return err
}
