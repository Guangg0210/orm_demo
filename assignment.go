package orm

/*
Assignable 标记接口
表达赋值语句，用于 UPDATE 和 UPSERT
*/
type Assignable interface {
	assign()
}

type Assignment struct {
	column string
	val    Expression
}

func (a Assignment) assign() {}

func Assign(column string, val any) Assignment {
	v, ok := val.(Expression)
	if !ok {
		v = value{val: val}
	}
	return Assignment{
		column: column,
		val:    v,
	}
}
