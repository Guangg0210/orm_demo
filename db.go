package orm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"go.uber.org/multierr"
	"log"
	"orm_demo/internal/errs"
	"orm_demo/internal/valuer"
	"orm_demo/model"
	"time"
)

type DBOption func(*DB)

type DB struct {
	core
	sqlDB *sql.DB
}

/*
Wait 等待数据库链接，
只用作测试。
*/
func (db *DB) Wait() error {
	err := db.sqlDB.Ping()
	for err == driver.ErrBadConn {
		log.Printf("等待数据库启动...")
		err = db.sqlDB.Ping()
		time.Sleep(time.Second)
	}
	return err
}

/*
Open 创建一个 DB 实例。
默认情况下，该 DB 将使用 MySQL 方言
如果使用了其他数据库，可以使用 DBWithDialect 制定方言
*/
func Open(driver string, dsn string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}

// OpenDB 方便测试
// 利用 mockDB
func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		core: core{
			r:          model.NewRegistry(),
			valCreator: valuer.NewUnsafeValue,
			dialect:    MySQL,
		},
		sqlDB: db,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

// DBWithDialect 指定方言
func DBWithDialect(dialect Dialect) DBOption {
	return func(db *DB) {
		db.dialect = dialect
	}
}

func DBWithRegistry(r model.Registry) DBOption {
	return func(db *DB) {
		db.r = r
	}
}

func DBUseReflectValuer() DBOption {
	return func(db *DB) {
		db.valCreator = valuer.NewReflectValue
	}
}

func DBWithMiddlewares(ms ...Middleware) DBOption {
	return func(db *DB) {
		db.ms = ms
	}
}

// MustNewDB 创建一个 DB，如果失败则会 panic
// 我个人不太喜欢这种
func MustNewDB(driver string, dsn string, opts ...DBOption) *DB {
	db, err := Open(driver, dsn, opts...)
	if err != nil {
		panic(err)
	}
	return db
}

// BeginTx 开启事物
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.sqlDB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{
		tx: tx,
	}, nil
}

/*
DoTx 将会开启事物 fn。 如果 fn 返回错误或者发生 panic，事物将会回滚，否则提交事物
*/
func (db *DB) DoTx(ctx context.Context, opts *sql.TxOptions,
	fn func(ctx context.Context, tx *Tx) error) (err error) {

	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	panicked := true
	defer func() {
		// 这里也可以用 recover()
		if panicked || err != nil {
			er := tx.Rollback()
			if er != nil {
				err = errs.NewErrFailToRollbackTx(err, er, panicked)
			}
			err = multierr.Combine(err, er)
		} else {
			err = multierr.Combine(err, tx.Commit())
		}
	}()

	err = fn(ctx, tx)
	panicked = false
	return
}

type txKey struct {
	//
}

// BeginTxV2 事务扩散的解决方案，
//func (db *DB) BeginTxV2(ctx context.Context,
//	opts *sql.TxOptions) (context.Context, *Tx, error) {
//	val := ctx.Value(txKey{})
//	if val != nil {
//		tx := val.(*Tx)
//		if !tx.done {
//			return ctx, tx, nil
//		}
//	}
//	tx, err := db.BeginTx(ctx, opts)
//	if err != nil {
//		return ctx, nil, err
//	}
//	ctx = context.WithValue(ctx, txKey{}, tx)
//	return ctx, tx, nil
//}

func (db *DB) Close() error {
	return db.sqlDB.Close()
}

func (db *DB) getCore() core {
	return db.core
}
func (db *DB) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.sqlDB.QueryContext(ctx, query, args...)
}

func (db *DB) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.sqlDB.ExecContext(ctx, query, args...)
}
