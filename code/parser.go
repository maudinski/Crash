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
    nodes := make([]node, 0)//IDEA append will eventually be a large bottle neck
    p.lx.flow()
    for t := range p.lx.stream {
        switch t.ttype {
        case "TYPE":
            n := p.getDeclarationN(t)
            nodes = append(nodes, n)
        case "PRINT":
            n := p.getPrintN()
            nodes = append(nodes, n)
        case "NAME":
            n := p.getReassignN(t)
            nodes = append(nodes, n)
        }
    }
    ast.statements = nodes
    return ast
}

func (p *parser) getReassignN(nameTok token) Reassign {
    _ = <- p.lx.stream
    valueTok := <- p.lx.stream
    return Reassign{Name{nameTok.value}, Value{valueTok.value}}
}

//absolutely no error checking
//tokens will/should be:
//(passed in){TTYPE, "int"}
//{"NAME", "x"}
//{"EQUAL", "="}
//{NUMBER, "5"}
//this assumes no operation
func (p *parser) getDeclarationN(t token) Declaration {
    nameTok := <- p.lx.stream
    eqTok := <-p.lx.stream
    if eqTok.ttype != "EQUAL"{//BUG it already isnt called this
        //throw an error
    }
    numTok := <-p.lx.stream
    questTok := <- p.lx.stream //check if its an operator token
    if questTok.ttype == "NEWLINE"{
        return Declaration{Name{nameTok.value}, Type{t.value}, Value{numTok.value}}
    }
    if questTok.ttype == "OPERATOR"{
        rightTok := <- p.lx.stream
        o := Operation{Operator{Value{questTok.value}}, Value{numTok.value}, Value{rightTok.value}}
        return Declaration{Name{nameTok.value}, Type{t.value}, o}
    }
    return Declaration{}
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

























