package orm

// RawExpr 代表一个原生表达式
// 意味着 ORM 不会对它进行任何处理
type RawExpr struct {
	raw  string
	args []any
}

func (r RawExpr) fieldName() string {
	return ""
}

func (r RawExpr) selectedAlias() string {
	return ""
}

func (r RawExpr) target() TableReference {
	return nil
}

func (r RawExpr) assign() {}

func (r RawExpr) expr() {}

func (r RawExpr) AsPredicate() Predicate {
	return Predicate{
		left: r,
	}
}

// Raw 创建一个 RawExpr
func Raw(expr string, args ...interface{}) RawExpr {
	return RawExpr{
		raw:  expr,
		args: args,
	}
}

/*
binaryExpr
作为 MathExpr 和 Predicate 的公共抽象，
代表二元操作符 a op b 的形态。
*/
type binaryExpr struct {
	left  Expression
	op    op
	right Expression
}

func (binaryExpr) expr() {}

/*
MathExpr 用于 UPDATE 语句
实现了 Expression 接口和 Assignable 接口，Column 本身是构建 MathExpr 的起点
*/
type MathExpr binaryExpr

func (MathExpr) expr() {}

func (m MathExpr) Add(val any) MathExpr {
	return MathExpr{
		left:  m,
		op:    opADD,
		right: valueOf(val),
	}
}

func (m MathExpr) Mutil(val any) MathExpr {
	return MathExpr{
		left:  m,
		op:    opMulti,
		right: valueOf(val),
	}
}

type SubQueryExpr struct {
	s SubQuery
	// 谓词 ALL、ANY 或者 SOME
	pred string
}

func (SubQueryExpr) expr() {}

func Any(sub SubQuery) SubQueryExpr {
	return SubQueryExpr{
		s:    sub,
		pred: "ANY",
	}
}

func All(sub SubQuery) SubQueryExpr {
	return SubQueryExpr{
		s:    sub,
		pred: "ALL",
	}
}
func Some(sub SubQuery) SubQueryExpr {
	return SubQueryExpr{
		s:    sub,
		pred: "SOME",
	}
}
