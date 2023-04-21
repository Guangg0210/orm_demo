package orm

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTx(t *testing.T) {
	db := memoryDB(t)
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	})
	assert.NoError(t, err)
	// 带事务的查询
	NewInserter[TestModel](tx)
	// 不带事务的查询
	NewInserter[TestModel](db)

}

/*
事务的使用方法1
type UserDAO struct {
	db *sql.DB
}

func (u *UserDAO) Begin() *UserTxDAO {

}
func (u *UserDAO) GetById() *UserTxDAO {

}

type UserTxDAO struct {
	tx *sql.Tx
	UserDAO
}

func (u *UserTxDAO) Commit() error {

}
func (u *UserTxDAO) Rollback() error {

}
*/

type user struct {
	id   uint64
	name string
}

// 处理方法2
type UserDAO struct {
	sess interface {
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	}
	user struct {
		id   uint64
		name string
	}
}

func (dao *UserDAO) GetById(ctx context.Context, id uint64) (*user, error) {
	_, err := dao.sess.QueryContext(ctx, "SELECT name FROM user WHERE id = ?", dao.user.id)
	if err != nil {
		return nil, err
	}
	// 处理结果集
	return &user{}, nil
}

func TestUserTxDAO(t *testing.T) {
	//db := memoryDB(t)
	//userDao := &UserDAO{
	//	sess: db.sqlDB,
	//}
	//
	//tx, err := db.BeginTx(context.Background(), &sql.TxOptions{
	//	Isolation: 0,
	//	ReadOnly:  false,
	//})
	//assert.NoError(t, err)
	//userTxDAO := &UserDAO{
	//	sess: tx.tx,
	//}

}
