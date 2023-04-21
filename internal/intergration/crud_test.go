package intergration

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	orm "orm_demo"
	"orm_demo/internal/test"
	"testing"
)

func TestMySQLCRUD(t *testing.T) {
	db, err := orm.Open("mysql", "root:Guan0210@tcp(localhost:13306)/integration_test")
	if err != nil {
		t.Fatal(err)
	}
	db.Wait()

	testCases := []struct {
		name     string
		i        *orm.Inserter[test.SimpleStruct]
		affected int64
		wantErr  error

		wantData *test.SimpleStruct
	}{
		{
			name:     "insert single",
			i:        orm.NewInserter[test.SimpleStruct](db).Values(test.NewSimpleStruct(12)),
			affected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.i.Exec(context.Background())
			affected, err := res.RowsAffected()
			if err != nil {
				return
			}
			assert.Equal(t, tc.affected, affected)
			id, err := res.LastInsertId()
			if err != nil {
				return
			}

			data, err := orm.NewSelector[test.SimpleStruct](db).Where(orm.C("Id").EQ(id)).Get(context.Background())
			require.NoError(t, err)
			assert.Equal(t, tc.wantData, data)
		})
	}
}
