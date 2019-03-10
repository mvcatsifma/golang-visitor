package main

import "fmt"

func main() {
	var links []string
	m1 := &Module{name: "1"}
	m1.modules = []*Module{
		{name: "3"},
	}
	root := &Module{
		name: "root",
		modules: []*Module{
			m1,
			{name: "2"},
		},
	}
	Inspect(root, func(module IModule) bool {
		switch m := module.(type) {
		case *Module:
			fmt.Printf("inspecting module: %v \n", m.name)
			links = append(links, m.getLinks())
		default:
			return false
		}
		return true
	})

	fmt.Println("Links....")
	for _, v := range links {
		fmt.Println(v)
	}
}

type Visitor interface {
	Visit(module IModule) (w Visitor)
}

type IModule interface {
	isModule()
	getLinks() string
}

type Module struct {
	name    string
	modules []*Module
}

func (this *Module) getLinks() string {
	return fmt.Sprintf("/api/%v", this.name)
}

func (*Module) isModule() {
}

type inspector func(IModule) bool

func (f inspector) Visit(module IModule) Visitor {
	fmt.Println("visit")
	if f(module) {
		return f
	}
	return nil
}

func Walk(v Visitor, module IModule) {
	fmt.Println("walk")
	if v = v.Visit(module); v == nil {
		return
	}

	switch n := module.(type) {
	case *Module:
		for _, m := range n.modules {
			Walk(v, m)
		}
	default:
		panic(fmt.Sprintf("ast.Walk: unexpected node type %T", n))
	}

	v.Visit(nil)
}

func Inspect(module IModule, f func(IModule) bool) {
	Walk(inspector(f), module)
}
