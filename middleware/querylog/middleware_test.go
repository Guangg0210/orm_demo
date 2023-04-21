package querylog

import (
	"context"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	orm "orm_demo"
	mock_demo "orm_demo/middleware/mock"
	"testing"
	"time"
)

func TestMiddleware_build(t *testing.T) {
	builder := MiddlewareBuilder{}
	builder.LogFunc(NewBuilder().logFunc).SlowQueryThreshold(10)

	db, err := orm.Open("sqlite3",
		"file:test.sqlDB?cache=shared&mode=memory",
		orm.DBWithMiddlewares(builder.Build(),
			func(next orm.HandlerFunc) orm.HandlerFunc {
				return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
					time.Sleep(time.Second)
					return next(ctx, qc)
				}
			}))
	if err != nil {
		t.Fatal(err)
	}

	_, err = orm.NewSelector[TestModel](db).Get(context.Background())
	assert.NotNil(t, err)
	if err != nil {
		return
	}
}

type TestModel struct {
}

// TestMiddleware_Mock 不向数据库发起查询
func TestMiddleware_Mock(t *testing.T) {
	builder := &MiddlewareBuilder{}
	_, err := orm.Open("sqlite3",
		"file:test.sqlDB?cache=shared&mode=memory",
		orm.DBWithMiddlewares(builder.Build()))
	if err != nil {
		t.Fatal(err)
	}

	doBusiness(context.WithValue(context.Background(), mock_demo.MockKey{}, mock_demo.Mock{
		Result: &orm.QueryResult{
			Result: "返回一些和用户相关的数据",
		},
	}))
}

func doBusiness(input ...any) {}
