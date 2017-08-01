package main

import (
    "strconv"
)

var presedi = map[string]int{
    "==": 0
    ">=": 0
    "<=": 0
    ">": 0
    "<": 0
    "+": 1
    "-": 1
    "*": 2
    "/": 2
}
// passed -1 to start
func (p *Parser) parseExpression(prevPresedence int) Expression {
    t := p.lx.next()
    expFunc := p.expObjFuncs[t.ttype]
    if expFunc == nil {
        p.errorTrashExpression() // TODO make this trash the line and put back any { and the \n after if it's there
        return Id{} // return anything
    }
    exp := expFunc(t) // gives the approriate Expression node
    if t.lx.peek().ttype == "OPERATOR" {
        t = p.lx.next()
        ie := InfixExpression{t: t, op: t.value, left: exp}
        ie.right = p.parseExpression()
        if t.op has higher presedence than prePResedence {

        }
        return ie
    } else {
        return exp
    }
}

func newId(t token) Expression {
    return Id{t, t.value}
}

// probably. ParseFunctionCall also calls parseExpression so interesting amount of
// recursion going on
func newCall(t token) Expression {
    return p.parseFunctionCall(t)
}

func newString(t token) Expression {
    return String{t, t.value}
}

func newInt(t token) Expression {
    i, _ := strconv.Atoi(t.value)
    return Int{t, i}
}

func newFloat(t token) Expression {
    t, _ :=strconv.ParseFloat(t.value, 64)
    return Float{t, f}
}










// needs to leave open breacket { for block parsing
func (p *Parser) errorTrashExpressions(t, fmt string. args ...interface{}) {

}

