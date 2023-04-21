package orm

import (
	"context"
)

// RawQuerier 原生查询器
type RawQuerier[T any] struct {
	core
	sql  string
	args []any
	sess session
}

func (r *RawQuerier[T]) Build() (*Query, error) {
	return &Query{
		SQL:  r.sql,
		Args: r.args,
	}, nil
}

func (r *RawQuerier[T]) Exec(ctx context.Context) Result {
	return exec(ctx, r.sess, r.core, &QueryContext{
		Builder: r,
		Type:    "RAW",
	})
}

// RawQuerier[TestModel]("SELECT * FROM xxx WHERE xxx).Get(ctx)
func RawQuery[T any](sess session, sql string, args ...any) *RawQuerier[T] {
	return &RawQuerier[T]{
		sql:  sql,
		args: args,
		sess: sess,
		core: sess.getCore(),
	}
}

func (r *RawQuerier[T]) Get(ctx context.Context) (*T, error) {
	res := get[T](ctx, r.core, r.sess, &QueryContext{
		Type:    "RAW",
		Builder: r,
	})
	if res.Result != nil {
		return res.Result.(*T), res.Err
	}
	return nil, res.Err
}

func (r *RawQuerier[T]) GetMulti(ctx context.Context) ([]*T, error) {
	panic("implement me")
}
