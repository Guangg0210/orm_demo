package errs

import (
	"errors"
	"fmt"
)

var (
	ErrInputNil                  = errors.New("orm:不支持 nil ")
	ErrPointerOnly               = errors.New("orm:不支持的类型,只支持指针")
	ErrEmptyTableName            = errors.New("orm:表名为空")
	ErrEmptyFieldName            = errors.New("orm:列名为空")
	ErrNonCorrectType            = errors.New("orm:非正确类型")
	ErrInvalidTagContent         = errors.New("orm:错误的标签设置")
	errUnsupportedExpressionType = errors.New("orm:不支持的表达式")
	errUnknownField              = errors.New("orm:未知字段")
	errUnknownColumn             = errors.New("orm:未知列")
	errUnsupportedSelectable     = errors.New("orm:不支持的目标列")
	errUnsupportedAssignableType = errors.New("orm:不支持的 Assignable 表达式")
	errInvalidAddress            = errors.New("orm:无效地址")

	ErrTooManyReturnedColumns = errors.New("orm:过多列")
	ErrNoRows                 = errors.New("orm:未找到数据")
	ErrInsertZeroRow          = errors.New("orm:插入 0 行")
	ErrNoUpdatedColumns       = errors.New("orm: 未指定更新的列")
)

// NewErrUnsupportedExpressionType 返回一个不支持该  expression 错误信息
// 这样的写法用户可以 var ok = errors.Is(NewErrUnsupportedExpressionType("abc"), errUnsupportedExpressionType)

// NewErrUnsupportedExpressionType 不支持的表达式
func NewErrUnsupportedExpressionType(exp any) error {
	return fmt.Errorf("%w %v", errUnsupportedExpressionType, exp)
}

/*
后续可以加入状态码方便排查错误，因此需要设计结构体
*/

type OrmErr struct {
	code string
	msg  string
}

func (o OrmErr) Error() string {
	return fmt.Sprintf("orm-%v: %v", o.code, o.msg)
}
func NewErrUnknownField(name string) error {
	return fmt.Errorf("%w %v", errUnknownField, name)
}

func NewErrInvalidTagContent(tag string) error {
	return fmt.Errorf("%w %s", ErrInvalidTagContent, tag)
}

func NewErrUnknownColumn(name string) error {
	return fmt.Errorf("%w %s", errUnknownColumn, name)
}
func NewErrUnsupportedSelectable(exp any) error {
	return fmt.Errorf("%w %s", errUnsupportedSelectable, exp)
}

func NewErrUnsupportedAssignableType(exp any) error {
	return fmt.Errorf("%w %v", errUnsupportedAssignableType, exp)
}
func NewErrInvalidAddress(field string) error {
	return fmt.Errorf("%w %v", errInvalidAddress, field)
}

func NewErrFailToRollbackTx(bizErr error, rbErr error, panicked bool) error {
	return fmt.Errorf("orm: 回滚事务失败, 业务错误 %w, 回滚错误 %s, panic: %t",
		bizErr, rbErr.Error(), panicked)
}
