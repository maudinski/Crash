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

//TODO make this a switch statment
func (p *parser) parse() *AbSynTree {
    ast := new(AbSynTree)
    nodes := make([]node, 0)
    p.lx.flow()
    for t := range p.lx.stream {
        if t.ttype == "TYPE" {
            n := p.getDeclarationN(t)
            nodes = append(nodes, n)
        } else if t.ttype == "PRINT" {
            n := p.getPrintN()
            nodes = append(nodes, n)
        }
    }
    ast.statements = nodes
    return ast
}

//absolutely no error checking
//tokens will/should be:
//(passed in){TTYPE, "int"}
//{"NAME", "x"}
//{"EQUAL", "="}
//{NUMBER, "5"}
func (p *parser) getDeclarationN(t token) Declaration {
    nameTok := <- p.lx.stream
    _ = <-p.lx.stream
    numTok := <-p.lx.stream
    return Declaration{Name{nameTok.value}, Type{t.value}, Value{numTok.value}}
}

func (p *parser) getPrintN() Print {
    _ = <- p.lx.stream
    nodes := make([]node, 0)
    done := false
    for tok := range p.lx.stream {
        switch tok.ttype {
        case ")":
            done = true
        case "NAME":
            nodes = append(nodes, Name{tok.value})
        }
        if done {
            break
        }
    }
    return newPrint(nodes...)
}

























