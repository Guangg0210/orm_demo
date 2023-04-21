package orm

type Column struct {
	table TableReference
	name  string
	alias string
}

func (c Column) fieldName() string {
	return c.name
}

func (c Column) selectedAlias() string {
	return c.alias
}

func (c Column) target() TableReference {
	return c.table
}

func (c Column) expr() {}

func (c Column) assign() {}

// As 结构化别名
func (c Column) As(alias string) Column {
	return Column{
		name:  c.name,
		alias: alias,
	}
}

type value struct {
	val any
}

func (c value) expr() {}

func valueOf(val any) value {
	return value{
		val: val,
	}
}

// C 结构化列。
func C(name string) Column {
	return Column{name: name}
}

// EQ =，结构化 = 符号。
// 例如 C("id").Eq(12)，得到 id!= 12。
func (c Column) EQ(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opEQ,
		right: exprOf(arg),
	}
}

// NEQ !=，结构化 != 符号。
// 例如 C("id").NotEq(12)，得到 id!= 12。
func (c Column) NEQ(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opNEQ,
		right: exprOf(arg),
	}
}

// LT <，结构化 < 符号。
// 例如 C("id").LT(12) 得到 id < 12。
func (c Column) LT(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opLT,
		right: exprOf(arg),
	}
}

// LTEQ <=，结构化 <= 符号。
// 例如 C("id").LTEQ(12)，得到 id <= 12。
func (c Column) LTEQ(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opLTEQ,
		right: exprOf(arg),
	}
}

// GT >，结构化 > 符号。
// 例如 C("id").GT(12)，得到 id > 12。
func (c Column) GT(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opGT,
		right: exprOf(arg),
	}
}

// GTEQ >=，结构化 >= 符号。
// 例如 C("id").GTEQ(12)，得到 id >= 12。
func (c Column) GTEQ(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opGTEQ,
		right: exprOf(arg),
	}
}

// Add +，结构化 + 符号。
// 例如 C("id").Add(12)，得到 id + 12。
func (c Column) Add(delta int) MathExpr {
	return MathExpr{
		left:  c,
		op:    opADD,
		right: exprOf(delta),
	}
}

// Multi *，结构化 * 符号。
// 例如 C("id").Multi(12)，得到 id * 12。
func (c Column) Multi(delta int) MathExpr {
	return MathExpr{
		left:  c,
		op:    opMulti,
		right: exprOf(delta),
	}
}

// In 结构化 IN
// In 有两种输入，一种是 IN 自查询，另一种是普通的值。
// 可以定义两个方法，In 和 InQuery，也可以定义一种
// 这里使用一个方法
func (c Column) In(vals ...any) Predicate {
	return Predicate{
		left:  c,
		op:    opIN,
		right: exprOf(vals),
	}
}

func (c Column) InQuery(sub SubQuery) Predicate {
	return Predicate{
		left:  c,
		op:    opIN,
		right: sub,
	}
}
