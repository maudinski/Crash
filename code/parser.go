package main

import (

)

type parser struct {
    lx *lexer
}

func newParser(lx *lexer) *parser {
    p := new(parser)
    p.lx = lx
    return p
}

func (p *parser) parse() *AbSynTree {
    ast := new(AbSynTree)
    nodes := make([]node, 0)
    p.lx.flow()
    for t := range p.lx.stream {
        if t.ttype == "TYPE" {
            n := p.getDeclaration(t)
            nodes = append(nodes, n)
        } else if t.ttype == "PRINT" {
            n := p.getPrint()
            nodes = append(nodes, n)
        }
    }
    ast.statements = nodes
    return ast
}

func (p *parser) getDeclaration(t token) Declaration {
    d := Declaration{}
    d.ttype = t.value
    t = <- p.lx.stream
    d.name = t.value
    t = <- p.lx.stream
    d.value = t.value
    return d
}

func (p *parser) getPrint() Print {
    _ = <- p.lx.stream
    t := <- p.lx.stream
    return Print{t.value}
}


























