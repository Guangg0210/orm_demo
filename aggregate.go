package orm

// Aggregate 代表聚合函数，例如 AVG, MAX, MIN 等
type Aggregate struct {
	table TableReference
	fn    string
	arg   string
	alias string
}

func (a Aggregate) fieldName() string {
	return a.arg
}
func (a Aggregate) selectedAlias() string {
	return a.alias
}
func (a Aggregate) target() TableReference {
	return a.table
}

func (a Aggregate) expr() {}

// Avg 结构化聚合函数 AVG
func Avg(c string) Aggregate {
	return Aggregate{
		fn:  "AVG",
		arg: c,
	}
}

// Max 结构化聚合函数 MAX
func Max(c string) Aggregate {
	return Aggregate{
		fn:  "MAX",
		arg: c,
	}
}

// Min 结构化聚合函数 MIN
func Min(c string) Aggregate {
	return Aggregate{
		fn:  "MIN",
		arg: c,
	}
}

// Count 结构化聚合函数 COUNT
func Count(c string) Aggregate {
	return Aggregate{
		fn:  "COUNT",
		arg: c,
	}
}

// Sum 结构化聚合函数 SUM
func Sum(c string) Aggregate {
	return Aggregate{
		fn:  "SUM",
		arg: c,
	}
}

// As 为结构化聚合函数创建别名
func (a Aggregate) As(alias string) Aggregate {
	return Aggregate{
		fn:    a.fn,
		arg:   a.arg,
		alias: alias,
	}
}

// EQ =
func (a Aggregate) EQ(arg any) Predicate {
	return Predicate{
		left:  a,
		op:    opEQ,
		right: exprOf(arg),
	}
}

// NEQ !=
func (a Aggregate) NEQ(arg any) Predicate {
	return Predicate{
		left:  a,
		op:    opNEQ,
		right: exprOf(arg),
	}
}

// LT <
func (a Aggregate) LT(arg any) Predicate {
	return Predicate{
		left:  a,
		op:    opLT,
		right: exprOf(arg),
	}
}

// LTEQ <=
func (a Aggregate) LTEQ(arg any) Predicate {
	return Predicate{
		left:  a,
		op:    opLTEQ,
		right: exprOf(arg),
	}
}

// GT >
func (a Aggregate) GT(arg any) Predicate {
	return Predicate{
		left:  a,
		op:    opGT,
		right: exprOf(arg),
	}
}

// GTEQ >=
func (a Aggregate) GTEQ(arg any) Predicate {
	return Predicate{
		left:  a,
		op:    opGTEQ,
		right: exprOf(arg),
	}
}
