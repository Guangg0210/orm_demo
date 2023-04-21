package orm

import (
	"context"
	"database/sql"
	"orm_demo/internal/valuer"
	"orm_demo/model"
)

type core struct {
	r          model.Registry
	valCreator valuer.Creator
	dialect    Dialect
	ms         []Middleware
}

func getHandler[T any](ctx context.Context, sess session, c core, qc *QueryContext) *QueryResult {
	q, err := qc.Builder.Build()
	if err != nil {
		return &QueryResult{Err: err}
	}

	rows, err := sess.queryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return &QueryResult{Err: err}
	}

	if !rows.Next() {
		return &QueryResult{
			Err: ErrNoRows,
		}
	}

	tp := new(T)
	meta, err := c.r.Get(tp)
	if err != nil {
		return &QueryResult{Err: err}
	}

	val := c.valCreator(tp, meta)
	err = val.SetColumns(rows)

	return &QueryResult{
		Result: tp,
		Err:    err,
	}
}

func get[T any](ctx context.Context, c core, sess session, qc *QueryContext) *QueryResult {
	var handler HandlerFunc = func(ctx context.Context, qc *QueryContext) *QueryResult {
		return getHandler[T](ctx, sess, c, qc)
	}
	ms := c.ms
	for i := len(ms) - 1; i >= 0; i-- {
		handler = ms[i](handler)
	}
	return handler(ctx, qc)
}

func exec(ctx context.Context, sess session, c core, qc *QueryContext) Result {
	var handler HandlerFunc = func(ctx context.Context, qc *QueryContext) *QueryResult {
		return execHandler(ctx, sess, c, qc)
	}
	ms := c.ms
	for i := len(ms) - 1; i >= 0; i-- {
		handler = ms[i](handler)
	}
	qr := handler(ctx, qc)
	var res sql.Result
	if qr.Result != nil {
		res = qr.Result.(sql.Result)
	}
	return Result{
		err: qr.Err,
		res: res,
	}
}

func execHandler(ctx context.Context, sess session, c core, qc *QueryContext) *QueryResult {
	q, err := qc.Builder.Build()
	if err != nil {
		return &QueryResult{
			Err: err,
		}
	}
	res, err := sess.execContext(ctx, q.SQL, q.Args...)
	return &QueryResult{Err: err, Result: res}
}
