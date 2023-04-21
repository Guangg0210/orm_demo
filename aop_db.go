//go:build AopDB

package orm

import (
	"context"
	"database/sql"
)

type AopDB struct {
	db *sql.DB
	ms []Middleware
}

type AopDBContext struct {
	query string
	args  []any
}

type Handler func(ctx *AopDBContext) *AopDBResult

type Middleware func(next Handler) Handler

type AopDBResult struct {
	row *sql.Rows
}

func (db *AopDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Rows {
	var handler Handler = func(actx *AopDBContext) *AopDBResult {
		row, _ := db.db.QueryContext(ctx, actx.query, actx.args...)
		return &AopDBResult{
			row: row,
		}
	}
	for i := len(db.ms) - 1; i >= 0; i-- {
		handler = db.ms[i](handler)
	}
	res := handler(&AopDBContext{})
	return res.row
}
