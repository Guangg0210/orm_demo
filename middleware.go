package orm

import (
	"context"
	"orm_demo/model"
)

/*
Gorm AOP 方案

Create 有四个，分成两对:
• BeforeSave 和 AfterSave
• BeforeCreate 和 AfterCreate

Update 也是四个 Hook，分成两对:
• BeforeSave 和 AfterSave
• BeforeUpdate 和 AfterUpdate

Delete 有两个 Hook
它们构成了一对: • BeforeDelete 和 AfterDelete

Query 只有一个 Hook，就是 AfterFind。

• 分查询类型:对增删改查有不同的 Hook
• 分时机:在查询执行前，或者在查询执行后。这种顺序是预定义好的
• 修改上下文:每一个 Hook 内部都是可以修改执行上下文的。例如可以利用这个特性实现一个简单的分库分表中间件
这种设计，优点在于用户用起来还是比较简单的，例如使用 AfterUpdate 的时 候，可以很清楚确定这个会在 Update 语句的时候被调用。
缺点也很明显:
• 缺乏扩展性，用户指定不了顺序
• 如果 GORM 要扩展支持别的接入点，例如 BeforeFind，需要修改源码

这里使用的解决方案
• 抽象出来一个 QueryContext，代表查询上下文
• 抽象出来一个 QueryResult，代表查询结果
• 抽象出来 Handler，代表在这个上下文里面做点什么事情
• 抽象出来 Middleware，连接不同的 Handler

这种设计的缺陷就是用户实现 Middleware 的时候，可能存在大量的类型 断言之类的东西，或者需要自己判断是什么查询。
*/

type QueryContext struct {
	// Type 用来区分 SELECT、UPDATE、DELETE、INSERT
	Type string
	// Builder 使用的时候，大多数情况下，需要转换到具体的类型才能篡改查询
	// 这种需求需要把能够篡改的字段设置成公共字段暴露给用户
	Builder QueryBuilder
	Model   *model.Model
}

type QueryResult struct {
	// SELECT 语句返回值是 T 或者 []T
	// UPDATE、DELETE、INSERT 返回值是 Result
	Result any
	Err    error
}

type HandlerFunc func(ctx context.Context, qc *QueryContext) *QueryResult

type Middleware func(next HandlerFunc) HandlerFunc
