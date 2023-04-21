package querylog

import (
	"context"
	"log"
	orm "orm_demo"
	"orm_demo/utils"
	"time"
)

type MiddlewareBuilder struct {
	// 慢查询的毫秒单位
	threshold int64

	logFunc func(sql string, args ...any)
}

func (m *MiddlewareBuilder) SlowQueryThreshold(threshold int64) *MiddlewareBuilder {
	m.threshold = threshold
	return m
}

func (m *MiddlewareBuilder) LogFunc(logFunc func(sql string, args ...any)) *MiddlewareBuilder {
	m.logFunc = logFunc
	return m
}

func NewBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(sql string, args ...any) {
			if args == nil {
				log.Printf("%s\n", sql)
			} else {
				log.Printf("SQL:%s,args:%s\n", sql, args)
			}
		},
	}
}

func (m *MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.HandlerFunc) orm.HandlerFunc {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			start := time.Now()
			q, err := qc.Builder.Build()
			if err != nil {
				return &orm.QueryResult{
					Err: err,
				}
			}

			defer func() {
				duration := time.Now().Sub(start)
				// 设置了慢查询阈值，并且触发了
				if m.threshold > 0 && duration.Milliseconds() > m.threshold {
					m.logFunc(utils.Red("threshold"))
				}
			}()

			m.logFunc(q.SQL, q.Args...)
			res := next(ctx, qc)

			return res
		}
	}
}
