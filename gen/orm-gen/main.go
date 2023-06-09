package main

import (
	"bytes"
	_ "embed"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
)

// Go 会读取 tpl.gohtml 里面的内容填充到 genOrm 里面
//
//go:embed template.gohtml
var genOrm string

func main() {

	// 用户必须输入一个 src. 限制文件
	// 然后我们会在同目录下生产代码
	//src := os.Args[1]
	// 源代码目录，也是目标文件目录
	//dir := filepath.Dir(src)
	// 输入的目录
	//srcFileName := filepath.Base(src)

	src := "gen/orm-gen/testdata/user.go"
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, src, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	tv := &SingleFileEntryVisitor{}
	ast.Walk(tv, f)
	file := tv.Get()

	// 拿到解析后生成模版
	tpl := template.New("orm-gen")
	tpl, err = tpl.Parse(genOrm)
	if err != nil {
		panic(err)
	}
	bs := &bytes.Buffer{}

	err = tpl.Execute(bs, file)
	if err != nil {
		panic(err)
	}
}
