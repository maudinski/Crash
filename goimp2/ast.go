package main

import ()
 
// something like that
type Ast struct {
	structs   []string // as in fuck this for now
	globals   []Declaration
	functions map[string]*Function
}

func newAst() *Ast {
	ast := new(Ast)
	ast.structs = make([]string, 0)
	ast.globals = make([]Declaration, 0)
	ast.functions = make(map[string]*Function, 0) // i think
	return ast
}

// as long as im aware that this could take up a lot of memory in practce, it's cool to
// do it like this
func (ast *Ast) String() string {
	str := "Structs:"
	for _, s := range ast.structs {
		str += "\n" + s
	}
	str += "\n Globals:"
	for _, d := range ast.globals {
		str += "\n" + d.String()
	}
	str += "\n Functions: "
	for _, f := range ast.functions {
		str += "\n\n" + f.String()
	}
	return str
}
