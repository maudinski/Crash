package main

import (
    "strconv"
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
        case Reassign:
            ast.execReassign(n)
        }
    }
}

//BUG doesnt check if the value already exists, will crash if it doesnt
//also doesnt check if the types are the same
func (ast *AbSynTree) execReassign(n node) {
    r := n.(Reassign)
    ast.vars[r.n.value][1] = r.v.value
}


func (ast *AbSynTree) execDeclaration(n node) {
    d, _ := n.(Declaration)
    var val string
    switch d.v.(type){
    case Operation:
        o := d.v.(Operation)
        val = ast.execOperation(o)
    case Value:
        v := d.v.(Value)
        val = v.value
    }
    ast.vars[d.n.value] = []string{d.t.value, val}
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


func (ast *AbSynTree) execOperation(o Operation) string {
    vLeft, _ := o.vLeft.(Value)
    vRight, _ := o.vRight.(Value)
    iLeft, _ := strconv.Atoi(vLeft.value)
    iRight, _ := strconv.Atoi(vRight.value)
    var snum string
    switch o.op.v.value{
    case "+":
        snum = strconv.Itoa(iLeft + iRight)
    }
    return snum
}



















