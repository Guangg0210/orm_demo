package valuer

import (
	"database/sql"
	"orm_demo/internal/errs"
	"orm_demo/model"

	"reflect"
)

// reflectValue 基于反射的 Value
type reflectValue struct {
	val  reflect.Value
	meta *model.Model
}

var _ Creator = NewReflectValue

// NewReflectValue 返回一个封装好的，基于反射实现的 Value
// 输入 val 必须是一个指向结构体实例的指针，而不能是任何其它类型
func NewReflectValue(val interface{}, meta *model.Model) Value {
	return reflectValue{
		val:  reflect.ValueOf(val).Elem(),
		meta: meta,
	}
}

func (r reflectValue) Field(name string) (any, error) {
	res := r.val.FieldByName(name)
	if res == (reflect.Value{}) {
		return nil, errs.NewErrUnknownField(name)
	}
	return res.Interface(), nil
}

func (r reflectValue) SetColumns(rows *sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	if len(columns) > len(r.meta.FieldMap) {
		return errs.ErrTooManyReturnedColumns
	}

	// colValues 和 colEleValues 实质上最终都指向同一个对象
	colValues := make([]interface{}, len(columns))
	colEleValues := make([]reflect.Value, len(columns))
	for i, col := range columns {
		fd, ok := r.meta.ColumnMap[col]
		if !ok {
			return errs.NewErrUnknownColumn(col)
		}
		val := reflect.New(fd.Type)
		// Scan 函数需要指针
		colValues[i] = val.Interface()
		colEleValues[i] = val.Elem()
	}
	if err = rows.Scan(colValues...); err != nil {
		return err
	}
	for i, col := range columns {
		cm := r.meta.ColumnMap[col]
		fd := r.val.Field(cm.Index)
		fd.Set(colEleValues[i])
	}
	return nil
}
