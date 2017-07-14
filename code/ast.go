package main

import (

)

type AbSynTree struct {
    statements []node
    variables []Declaration
}

func (ast *AbSynTree) interpret() {
    for _, n := range ast.statements {
        ast.execute(n)
    }
}

func (ast *AbSynTree) execute(n node) interface{} {
    switch n.(type){
    case Declaration:
        d := n.(Declaration)
        ast.variables = append(ast.variables, d)//maybe type assertion
    case Print:
        p := n.(Print)
        p.exec()
    }
}
























