package orm

type TableReference interface {
	tableAlias() string
}

type Table struct {
	entity any
	alias  string
}

func TableOf(entity any) Table {
	return Table{
		entity: entity,
	}
}

func (t Table) C(name string) Column {
	return Column{
		name:  name,
		table: t,
	}
}

func (t Table) tableAlias() string {
	return t.alias
}

func (t Table) As(alias string) Table {
	return Table{
		entity: t.entity,
		alias:  alias,
	}
}

func (t Table) Join(target TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		right: target,
		typ:   "JOIN",
	}
}

func (t Table) LeftJoin(target TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		right: target,
		typ:   "LEFT JOIN",
	}
}

func (t Table) RightJoin(target TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		right: target,
		typ:   "RIGHT JOIN",
	}
}

type JoinBuilder struct {
	left  TableReference
	right TableReference
	typ   string
}

var _ TableReference = Join{}

/*
JOIN 语法有两种形态:
• JOIN ... ON
• JOIN ... USING:USING 后面使用的是列名
JOIN 本身有:
• INNER JOIN， JOIN
• LEFT JOIN，RIGHT JOIN
所以:
• SELECT 的 FROM 后面可以是一个普通的表
• 也可以是一个 JOIN 查询
*/

type Join struct {
	left  TableReference
	right TableReference
	typ   string
	on    []Predicate
	using []string
}

func (j Join) Join(target TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  j,
		right: target,
		typ:   "JOIN",
	}
}

func (j Join) LeftJoin(target TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  j,
		right: target,
		typ:   "LEFT JOIN",
	}
}

func (j Join) RightJoin(target TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  j,
		right: target,
		typ:   "RIGHT JOIN",
	}
}

func (j Join) tableAlias() string {
	return ""
}

func (j *JoinBuilder) On(ps ...Predicate) Join {
	return Join{
		left:  j.left,
		right: j.right,
		on:    ps,
		typ:   j.typ,
	}
}

func (j *JoinBuilder) Using(cs ...string) Join {
	return Join{
		left:  j.left,
		right: j.right,
		using: cs,
		typ:   j.typ,
	}
}

type SubQuery struct {
	// 使用 QueryBuilder 仅仅是为了让 SubQuery 可以是非泛型的。
	s       QueryBuilder
	columns []Selectable
	alias   string
	table   TableReference
}

func (s SubQuery) expr() {}

func (s SubQuery) tableAlias() string {
	return s.alias
}

func (s SubQuery) Join(target TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  s,
		right: target,
		typ:   "JOIN",
	}
}

func (s SubQuery) LeftJoin(target TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  s,
		right: target,
		typ:   "LEFT JOIN",
	}
}

func (s SubQuery) RightJoin(target TableReference) *JoinBuilder {
	return &JoinBuilder{
		left:  s,
		right: target,
		typ:   "RIGHT JOIN",
	}
}

func (s SubQuery) C(name string) Column {
	return Column{
		table: s,
		name:  name,
	}
}
