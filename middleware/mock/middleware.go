package mock_demo

import (
	"context"
	orm "orm_demo"
	"time"
)

type MiddlewareBuilder struct{}

func (m MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.HandlerFunc) orm.HandlerFunc {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			val := ctx.Value(MockKey{})

			// 如果用户设置了 mock ，这个中间件不会发起真的查询
			if val != nil {
				mock := val.(*Mock)
				if mock.Sleep > 0 {
					time.Sleep(mock.Sleep)
				}
				return mock.Result
			}
			return next(ctx, qc)
		}
	}
}

type Mock struct {
	Sleep  time.Duration
	Result *orm.QueryResult
}

type MockKey struct {
}
