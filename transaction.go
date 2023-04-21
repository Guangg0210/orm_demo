package orm

import (
	"context"
	"database/sql"
)

/*
Gorm 事物设计
DB 本身也可以被看做是事务
• 普通的事务开启、提交和回滚功能
• 额外实现了一个 SavePoint 的功能
• 事务闭包 API

如果上游的方法开启了事务，那么下游的所有方法也会使用这个事务，否则
• 下游可以开一个新事务
• 也可以无事务运行
• 还可以报错

• 事务扩散
	其实本质就是上下文里面有事务就用事务，没有事务就开新事务。
Go 里面要解决的话只能依赖于 context.Context，基本上在别的语言里面用 thread-local 解决
的，到 Go 里面都是用 context.Context

• 事务扩散中，如果没有开启事务应该怎么办?
	看你的业务，你可以选择报错，可以选择开启新事务，也可以无事务运行

• 事务重复提交会怎样?
	在 ORM 层面上，有些 ORM 会维护一个标记位，标记一个事务有没有被提交。即便没有这个标记位，数据库也会返回错误。

• Go 里面实现一个事务闭包要考虑一些什么问题?如何实现?
	主要是考虑 panic 的问题，而后要在panic 的时候，以及业务代码返回 error 的时候，回滚事务
*/

type Tx struct {
	core
	tx *sql.Tx
}

func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

/*
RollbackIfNotCommit
此时事务已经被提交，或者 被回滚掉了，那么就会得到 sql.ErrTxDone 错误， 这时候我们忽略这个错误就可以。
*/
func (tx *Tx) RollbackIfNotCommit() error {
	err := tx.tx.Rollback()
	if err == sql.ErrTxDone {
		return nil
	}
	return err
}
func (tx *Tx) getCore() core {
	return tx.core
}
func (tx *Tx) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return tx.tx.QueryContext(ctx, query, args...)
}

func (tx *Tx) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return tx.tx.ExecContext(ctx, query, args...)
}

type session interface {
	getCore() core
	queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	execContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}
