package filter

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

type User struct {
	Id    string `bun:",pk"`
	Name  string
	Phone string
	Email string
}

func newDB(t *testing.T) *bun.DB {
	sqlite, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}
	sqlite.SetMaxOpenConns(1)

	db := bun.NewDB(sqlite, sqlitedialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true), bundebug.FromEnv("BUNDEBUG"),
	))

	return db
}

func resetDB(ctx context.Context, db *bun.DB) error {
	if err := db.ResetModel(ctx, (*User)(nil)); err != nil {
		return err
	}

	seed := []User{
		{Id: "_id", Name: "google", Phone: "123456789", Email: "example@domain.com"},
		{Id: "_id", Name: "google", Phone: "123456789", Email: "example@domain.com"},
		{Id: "_id", Name: "google", Phone: "123456789", Email: "example@domain.com"},
		{Id: "_id", Name: "google", Phone: "123456789", Email: "example@domain.com"},
	}

	for i := range seed {
		seed[i].Id += fmt.Sprintf("_%d", i+1)
		seed[i].Name += fmt.Sprintf("_%d", i+1)
		seed[i].Phone += fmt.Sprintf("_%d", i+1)
		seed[i].Email += fmt.Sprintf("_%d", i+1)
	}

	if _, err := db.NewInsert().Model(&seed).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func TestALL(t *testing.T) {
	var (
		ctx = context.Background()
		db  = newDB(t)
		err = resetDB(ctx, db)
	)
	assert.NoError(t, err)

	var users []User
	q := db.NewSelect().Model(&users)
	{
		Limit(q, 10) // limit, order
		OrderByAsc(q, "id")
		OrderByDesc(q, "id")
		OrderBy(q, "id", "asc")
	}

	{
		Where(q, Equal("email", "example@domain.com"))
		Where(q, NotEqual("email", "example@domain.com"))

		Where(q, LessThan("email", 5))
		Where(q, LessThanOrEqual("email", 5))
		Where(q, GreaterThan("email", 5))
		Where(q, GreaterThanOrEqual("email", 5))

		OrWhere(q, Contains("name", "google"))
		OrWhere(q, NotContains("name", "google"))

		OrWhere(q, StartsWith("name", "goo"))
		OrWhere(q, NotStartsWith("name", "goo"))
		Where(q, EndsWith("email", "@domain.com"))
		Where(q, NotEndsWith("email", "@domain.com"))

		Where(q, In("phone", []string{"in_1", "in_2"}))
		OrWhere(q, In("phone", []string{"in_1", "in_2"}))
		Where(q, NotIn("phone", []string{"not_in_1", "not_in_2"}))
	}

	err = q.Scan(ctx)
	assert.Equal(t, nil, err)

	t.Log(users)
}
