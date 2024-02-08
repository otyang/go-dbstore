package dbstore

import "errors"

func NonNegativeLimit(limit int) int {
	if limit < 0 {
		return 0
	}
	return limit
}

type PaginationParams struct {
	Limit             int
	DirectionNextPage bool
	CursorColumn      string
	CursorValue       string
}

type PaginationOption func(o *PaginationParams) error

func WithCursor(limit int, directionNextPage bool, cursorColumn string, cursorValue string) PaginationOption {
	return func(o *PaginationParams) error {
		if limit < 0 {
			limit = 0
		}

		if cursorColumn == "" {
			return errors.New("cursor column not specified: must be specified")
		}

		o.Limit = limit
		o.DirectionNextPage = directionNextPage
		o.CursorColumn = cursorColumn
		o.CursorValue = cursorValue
		return nil
	}
}
