package obun

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	dbstore "github.com/otyang/go-dbstore"
	"github.com/otyang/go-dbstore/filter"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

// Use constants for configuration
const (
	testDBDriver = dbstore.DriverSqlite
	testDSN      = "file::memory:?cache=shared"
)

var seed = []Book{
	{Id: "1", Title: "Title 1"},
	{Id: "2", Title: "Title 2"},
	{Id: "3", Title: "Title 3"},
	{Id: "4", Title: "Title 4"},
}

type Book struct {
	Id    string `bun:",pk"`
	Title string `bun:",notnull"`
}

func setUpMigrateAndTearDown(t *testing.T, modelsPtr ...any) (context.Context, *bun.DB, *Repository, func()) {
	ctx := context.TODO()

	// connect
	db, err := dbstore.NewDBConnection(testDBDriver, testDSN, 1, true)
	assert.NoError(t, err)

	// migrate
	for _, model := range modelsPtr {
		_, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx)
		assert.NoError(t, err)
	}

	// tearDown
	teardownFunc := func() {
		for _, model := range modelsPtr {
			_, err := db.NewDropTable().Model(model).Exec(ctx)
			assert.NoError(t, err)
		}
	}

	return ctx, db, NewRepository(db), teardownFunc
}

func TestRepository_Create(t *testing.T) {
	ctx, _, repo, tearDown := setUpMigrateAndTearDown(t, (*Book)(nil))
	defer tearDown()

	data := Book{
		Id:    "_1234asdf",
		Title: "the unknown",
	}

	err := repo.Create(ctx, &data, false)
	assert.NoError(t, err)

	// re-inserting the same data should create an error
	// since the primary key already exists
	err = repo.Create(ctx, &data, false)
	assert.Error(t, err)

	// ignore duplicates
	err = repo.Create(ctx, &data, true)
	assert.NoError(t, err)
}

func TestRepository_CreateBulk(t *testing.T) {
	ctx, _, repo, tearDown := setUpMigrateAndTearDown(t, (*Book)(nil))
	defer tearDown()

	err := repo.CreateBulk(ctx, &seed, false)
	assert.NoError(t, err)

	err = repo.CreateBulk(ctx, &seed, false)
	assert.Error(t, err)

	// ignore duplicates
	err = repo.CreateBulk(ctx, &seed, true)
	assert.NoError(t, err)
}

func TestRepository_FindOneByPK(t *testing.T) {
	var (
		ctx, _, repo, tearDown = setUpMigrateAndTearDown(t, (*Book)(nil))
		err                    = repo.CreateBulk(ctx, &seed, false)
	)

	defer tearDown()
	assert.NoError(t, err)

	data := Book{Id: "1"}
	err = repo.FindOneByPK(ctx, &data)

	assert.NoError(t, err)
	assert.Equal(t, seed[0], data)
}

func TestRepository_FindOneWhere(t *testing.T) {
	var (
		ctx, _, repo, tearDown = setUpMigrateAndTearDown(t, (*Book)(nil))
		err                    = repo.CreateBulk(ctx, &seed, false)
		bookFromDB             Book
	)

	defer tearDown()
	assert.NoError(t, err)

	err = repo.FindOneWhere(ctx, &bookFromDB)
	assert.NoError(t, err)
	assert.NotEmpty(t, bookFromDB.Id)

	err = repo.FindOneWhere(ctx, &bookFromDB, func(q *bun.SelectQuery) *bun.SelectQuery {
		filter.Where(q, filter.Equal("id", "2"))
		return q
	})
	assert.NoError(t, err)
	assert.Equal(t, seed[1], bookFromDB)
}

func TestRepository_FindManyWhere(t *testing.T) {
	var (
		ctx, _, repo, tearDown = setUpMigrateAndTearDown(t, (*Book)(nil))
		err                    = repo.CreateBulk(ctx, &seed, false)
		Books                  []Book
	)

	defer tearDown()
	assert.NoError(t, err)

	t.Run("FindManyWhere without select criterias", func(t *testing.T) {
		err := repo.FindManyWhere(ctx, &Books, nil)
		assert.NoError(t, err)
		assert.Equal(t, seed, Books)
		assert.Equal(t, len(seed), len(Books))
	})

	t.Run("FindManyWhere with select criterias", func(t *testing.T) {
		var got []Book
		err := repo.FindManyWhere(ctx, &got, nil, func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("id >= ?", 2)
		})
		assert.NoError(t, err)
		assert.Equal(t, 3, len(got))
	})
}

func TestRepository_UpdateByPK_oneAndMany(t *testing.T) {
	var (
		ctx, _, repo, tearDown = setUpMigrateAndTearDown(t, (*Book)(nil))
		err                    = repo.CreateBulk(ctx, &seed, false)
	)

	defer tearDown()
	assert.NoError(t, err)

	t.Run("UpdateOne By PK", func(t *testing.T) {
		want := seed[0]
		want.Title = "Updated Title 1..."
		err := repo.UpdateOneByPK(ctx, &want)
		assert.NoError(t, err)

		got := Book{Id: "1"}
		err = repo.FindOneByPK(ctx, &got)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Title 1...", got.Title)
	})

	t.Run("Update Many By PK", func(t *testing.T) {
		updatedBooks := seed
		updatedBooks[2].Title = "bulk update 3"
		updatedBooks[3].Title = "bulk update 4"
		err := repo.UpdateManyByPK(ctx, &updatedBooks)
		assert.NoError(t, err)

		got := Book{Id: "3"}
		err = repo.FindOneByPK(ctx, &got)
		assert.NoError(t, err)
		assert.Equal(t, "bulk update 3", got.Title)
	})

	t.Run("UpdateOneWhere", func(t *testing.T) {
		want := seed[0]
		want.Title = "one where"
		err = repo.UpdateOneWhere(ctx, &want, func(q *bun.UpdateQuery) *bun.UpdateQuery {
			return q.Where("id = ?", seed[0].Id)
		})

		got := Book{Id: "1"}
		repo.FindOneByPK(ctx, &got)
		assert.NoError(t, err)
		assert.Equal(t, "one where", got.Title)
	})
}

func TestRepository_Upsert(t *testing.T) {
	var (
		ctx, _, repo, tearDown = setUpMigrateAndTearDown(t, (*Book)(nil))
		err                    = repo.CreateBulk(ctx, &seed, false)
	)

	defer tearDown()
	assert.NoError(t, err)

	upsertedBooks := seed
	upsertedBooks[3].Title = "bulk update 4 9"

	err = repo.Upsert(ctx, &upsertedBooks)
	assert.NoError(t, err)

	var gotListOfUpsertedBooks []Book
	err = repo.FindManyWhere(ctx, &gotListOfUpsertedBooks, nil)
	assert.NoError(t, err)
	assert.Equal(t, seed[3], gotListOfUpsertedBooks[3])
}

func TestRepository_DeleteByPK(t *testing.T) {
	var (
		ctx, _, repo, tearDown = setUpMigrateAndTearDown(t, (*Book)(nil))
		err                    = repo.CreateBulk(ctx, &seed, false)
	)

	defer tearDown()
	assert.NoError(t, err)

	t.Run("DeleteByPK", func(t *testing.T) {
		err = repo.DeleteByPK(ctx, &seed[0])
		assert.NoError(t, err)

		err = repo.FindOneByPK(ctx, &seed[0])
		assert.Equal(t, sql.ErrNoRows.Error(), err.Error())
	})

	t.Run("DeleteByPK  many", func(t *testing.T) {
		err = repo.DeleteByPK(ctx, &[]Book{seed[1], seed[2]})
		assert.NoError(t, err)

		err = repo.FindOneByPK(ctx, &seed[1])
		assert.Equal(t, sql.ErrNoRows.Error(), err.Error())

		err = repo.FindOneByPK(ctx, &seed[2])
		assert.Equal(t, sql.ErrNoRows.Error(), err.Error())
	})
}

func TestRepository_DeleteWhere(t *testing.T) {
	var (
		ctx, _, repo, tearDown = setUpMigrateAndTearDown(t, (*Book)(nil))
		err                    = repo.CreateBulk(ctx, &seed, false)
	)

	defer tearDown()
	assert.NoError(t, err)

	t.Run("DeleteWhere", func(t *testing.T) {
		err = repo.DeleteWhere(ctx, (*Book)(nil), func(q *bun.DeleteQuery) *bun.DeleteQuery {
			filter.Where(q, filter.Equal("id", 1))
			return q
		})

		assert.NoError(t, err)

		err = repo.FindOneByPK(ctx, &seed[0])
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows.Error(), err.Error())
	})
}

func TestRepository_Transaction(t *testing.T) {
	ctx, _, repo, tearDown := setUpMigrateAndTearDown(t, (*Book)(nil))
	defer tearDown()

	t.Run("transactions: no errors", func(t *testing.T) {
		err := repo.Transaction(ctx, func(ctx context.Context, tx bun.Tx) error {
			if err := repo.NewWithTx(tx).Create(ctx, &seed[0], false); err != nil {
				return err
			}
			if err := repo.NewWithTx(tx).Create(ctx, &seed[1], false); err != nil {
				return err
			}
			return repo.NewWithTx(tx).Create(ctx, &seed[2], false)
		},
		)
		assert.NoError(t, err)
	})

	t.Run("transactions: with deliberate error to abort transactions", func(t *testing.T) {
		err := repo.Transaction(ctx, func(ctx context.Context, tx bun.Tx) error {
			err := repo.NewWithTx(tx).Create(ctx, &seed[3], false)
			if err != nil {
				return err
			}

			assert.NoError(t, err)
			return errors.New("deliberate-wrong-data")
		},
		)
		assert.Error(t, err)
	})
}
