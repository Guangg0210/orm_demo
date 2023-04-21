package orm

import (
	"orm_demo/internal/errs"
)

var (
	MySQL   Dialect = &mysqlDialect{}
	SQLite3 Dialect = &sqlite3Dialect{}
)

type Dialect interface {
	quoter() byte
	buildUpsert(builder *builder, odk *Upsert) error
}

// standardSQL
type standardSQL struct {
}

type mysqlDialect struct {
	standardSQL
}

func (dialect *mysqlDialect) quoter() byte {
	return '`'
}

func (dialect *mysqlDialect) buildUpsert(b *builder, odk *Upsert) error {
	if len(odk.assigns) == 0 {
		return errs.ErrNoUpdatedColumns
	}

	b.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	for idx, a := range odk.assigns {
		if idx > 0 {
			b.sb.WriteByte(',')
		}
		switch assign := a.(type) {
		case Assignment:
			err := b.buildColumn(nil, assign.column)
			if err != nil {
				return err
			}
			b.sb.WriteString("=")
			return b.buildExpression(assign.val)
		case Column:
			colName, err := b.colName(assign.table, assign.name)
			if err != nil {
				return err
			}
			b.quote(colName)
			b.sb.WriteString("=VALUES(")
			b.quote(colName)
			b.sb.WriteByte(')')
		default:
			return errs.NewErrUnsupportedAssignableType(a)
		}
	}
	return nil
}

type sqlite3Dialect struct {
	standardSQL
}

func (s sqlite3Dialect) quoter() byte {
	return '`'
}

// buildDuplicateKey 构建 UPSERT ON CONFLICT 部分的代码，
// 参考资料 https://www.sqlite.org/lang_UPSERT.html
// TODO 没有完善
func (s sqlite3Dialect) buildUpsert(b *builder, odk *Upsert) error {
	if len(odk.assigns) == 0 {
		return errs.ErrNoUpdatedColumns
	}
	b.sb.WriteString("ON CONFLICT ")
	if len(odk.conflictColumns) > 0 {
		b.sb.WriteByte('(')
		for i, col := range odk.conflictColumns {
			if i > 0 {
				b.sb.WriteByte(',')
			}
			err := b.buildColumn(nil, col)
			if err != nil {
				return err
			}
		}
		b.sb.WriteByte(')')
		b.sb.WriteByte(' ')
	}
	b.sb.WriteString("DO")
	b.sb.WriteString(" UPDATE SET ")
	for idx, a := range odk.assigns {
		if idx > 0 {
			b.sb.WriteByte(',')
		}
		switch assign := a.(type) {
		case Assignment:
			err := b.buildColumn(nil, assign.column)
			if err != nil {
				return err
			}
			b.sb.WriteString("=")
			return b.buildExpression(assign.val)
		case Column:
			colName, err := b.colName(assign.table, assign.name)
			if err != nil {
				return err
			}
			b.quote(colName)
			b.sb.WriteString("=exclude.")
			b.quote(colName)
		default:
			return errs.NewErrUnsupportedAssignableType(assign)
		}
	}
	return nil
}
