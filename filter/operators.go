package filter

import "strings"

// shortcut alias. i.e alternative to longer names
var (
	Eq        = Equal
	NEq       = NotEqual
	Lt        = LessThan
	Lte       = LessThanOrEqual
	Gt        = GreaterThan
	Gte       = GreaterThanOrEqual
	Starts    = StartsWith
	NotStarts = NotStartsWith
	Ends      = EndsWith
	NotEnds   = NotEndsWith
)

type sqlWhere struct {
	stmt                string
	columnName          string
	columnValue         any
	sql_IsNullQueryType bool
}

func newSqlWhereStmt(skipStatement bool, stmt string, columnName string, columnValue any) *sqlWhere {
	if skipStatement {
		return nil
	}
	return &sqlWhere{
		stmt:        stmt,
		columnName:  columnName,
		columnValue: columnValue,
	}
}

func empty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func (n *sqlWhere) isANullQueryType() bool {
	return n.sql_IsNullQueryType
}

func Equal(columnName string, value any) *sqlWhere {
	return newSqlWhereStmt(value == nil, "? = ?", columnName, value)
}

func NotEqual(columnName string, value any) *sqlWhere {
	return newSqlWhereStmt(value == nil, "? != ?", columnName, value)
}

func LessThan(columnName string, value any) *sqlWhere {
	return newSqlWhereStmt(value == nil, "? < ?", columnName, value)
}

func LessThanOrEqual(columnName string, value any) *sqlWhere {
	return newSqlWhereStmt(value == nil, "? <= ?", columnName, value)
}

func GreaterThan(columnName string, value any) *sqlWhere {
	return newSqlWhereStmt(value == nil, "? > ?", columnName, value)
}

func GreaterThanOrEqual(columnName string, value any) *sqlWhere {
	return newSqlWhereStmt(value == nil, "? >= ?", columnName, value)
}

func Contains(columnName string, value string) *sqlWhere {
	return newSqlWhereStmt(empty(value), "lower(?) LIKE ?", columnName, "%"+value+"%")
}

func NotContains(columnName string, value string) *sqlWhere {
	return newSqlWhereStmt(empty(value), "lower(?) NOT LIKE ?", columnName, "%"+value+"%")
}

func StartsWith(columnName string, value string) *sqlWhere {
	return newSqlWhereStmt(empty(value), "lower(?) LIKE ?", columnName, value+"%")
}

func NotStartsWith(columnName string, value string) *sqlWhere {
	return newSqlWhereStmt(empty(value), "lower(?) NOT LIKE ?", columnName, value+"%")
}

func EndsWith(columnName string, value string) *sqlWhere {
	return newSqlWhereStmt(empty(value), "lower(?) LIKE ?", columnName, "%"+value)
}

func NotEndsWith(columnName string, value string) *sqlWhere {
	return newSqlWhereStmt(empty(value), "lower(?) NOT LIKE ?", columnName, "%"+value)
}

func In[T any](columnName string, value []T) *sqlWhere {
	return newSqlWhereStmt(len(value) == 0, "? IN (?)", columnName, value)
}

func NotIn[T any](columnName string, value []T) *sqlWhere {
	return newSqlWhereStmt(len(value) == 0, "? NOT IN (?)", columnName, value)
}

func IsNull(columnName string) *sqlWhere {
	q := newSqlWhereStmt(false, "? IS NULL", columnName, nil)
	q.sql_IsNullQueryType = true
	return q
}

func IsNotNull(columnName string) *sqlWhere {
	q := newSqlWhereStmt(false, "? IS NOT NULL", columnName, nil)
	q.sql_IsNullQueryType = true
	return q
}
