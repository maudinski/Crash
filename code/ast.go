package main

import (

)

type AbSynTree struct {
    statements []node
    vars map[string][]string//name: type, value
}

func (ast *AbSynTree) interpret() {
    ast.vars = make(map[string][]string)
    for _, n := range ast.statements {
        switch n.(type){
        case Declaration:
            ast.execDeclaration(n)
        case Print:
            ast.execPrint(n)
        }
    }
}

func (ast *AbSynTree) execDeclaration(n node) {
    d, _ := n.(Declaration)
    ast.vars[d.n.value] = []string{d.t.value, d.v.value}
}

//BUG if variable doesnt exist
func (ast *AbSynTree) execPrint(n node) {
    p, _ := n.(Print)
    var prints []string
    for _, n := range p.params {
        switch n.(type){
        case Name:
            name := n.(Name)
            prints = append(prints, name.value)
        }
    }
    if len(prints) == 0 { return }
    print(ast.vars[prints[0]][1])
    for i := 1; i < len(prints); i++ {
        print(" ")
        print(ast.vars[prints[i]][1])
    }
    print("\n")
}






















