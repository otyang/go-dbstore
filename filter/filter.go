package filter

import (
	"fmt"
	"strings"

	"github.com/uptrace/bun"
)

func Limit(q *bun.SelectQuery, limit int) {
	q.Limit(limit)
}

func OrderBy(q *bun.SelectQuery, column, direction string) {
	if strings.TrimSpace(direction) == "" {
		return
	}
	value := strings.ToLower(strings.TrimSpace(direction))

	switch value {
	case "asc":
		q.OrderExpr(column + " asc")
	case "desc":
		q.OrderExpr(column + " desc")
	default:
		panic("invalid order direction: should be 'asc' or 'desc'")
	}
}

func OrderByAsc(q *bun.SelectQuery, column string) {
	OrderBy(q, column, "asc")
}

func OrderByDesc(q *bun.SelectQuery, column string) {
	OrderBy(q, column, "desc")
}

type allBunQueryType interface {
	*bun.SelectQuery | *bun.UpdateQuery | *bun.DeleteQuery
}

func Where[T allBunQueryType](bunQ T, opt *sqlWhere) {
	if opt == nil {
		return
	}

	var (
		sqlQuery    = opt.stmt
		column      = bun.Ident(opt.columnName)
		columnValue = opt.columnValue
		isNullQuery = opt.isANullQueryType()
	)

	if !isNullQuery {
		switch q := any(bunQ).(type) {
		case *bun.SelectQuery:
			q.Where(sqlQuery, column, columnValue)
		case *bun.UpdateQuery:
			q.Where(sqlQuery, column, columnValue)
		case *bun.DeleteQuery:
			q.Where(sqlQuery, column, columnValue)
		default:
			fmt.Println("unsupported type: and where only works with Select, Update & Delete Query")
			panic("unsupported type: and where only works with Select, Update & Delete Query")
		}
	}

	if isNullQuery {
		switch q := any(bunQ).(type) {
		case *bun.SelectQuery:
			q.Where(sqlQuery, column)
		case *bun.UpdateQuery:
			q.Where(sqlQuery, column)
		case *bun.DeleteQuery:
			q.Where(sqlQuery, column)
		default:
			fmt.Println("unsupported type: and where only works with Select, Update & Delete Query")
			panic("unsupported type: and where only works with Select, Update & Delete Query")
		}
	}
}

func OrWhere[T allBunQueryType](bunQ T, opt *sqlWhere) {
	if opt == nil {
		return
	}

	var (
		sqlQuery    = opt.stmt
		column      = bun.Ident(opt.columnName)
		columnValue = opt.columnValue
		isNullQuery = opt.isANullQueryType()
	)

	if !isNullQuery {
		switch q := any(bunQ).(type) {
		case *bun.SelectQuery:
			q.WhereOr(sqlQuery, column, columnValue)
		case *bun.UpdateQuery:
			q.WhereOr(sqlQuery, column, columnValue)
		case *bun.DeleteQuery:
			q.WhereOr(sqlQuery, column, columnValue)
		default:
			fmt.Println("unsupported type: or where only works with Select, Update & Delete Query")
			panic("unsupported type: or where only works with Select, Update & Delete Query")
		}
	}

	if isNullQuery {
		switch q := any(bunQ).(type) {
		case *bun.SelectQuery:
			q.WhereOr(sqlQuery, column)
		case *bun.UpdateQuery:
			q.WhereOr(sqlQuery, column)
		case *bun.DeleteQuery:
			q.WhereOr(sqlQuery, column)
		default:
			fmt.Println("unsupported type: or where only works with Select, Update & Delete Query")
			panic("unsupported type: or where only works with Select, Update & Delete Query")
		}
	}
}
