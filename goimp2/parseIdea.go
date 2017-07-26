// this is shittily done. get it working then redo
package main

import (
	"fmt"
	"os"
)

// something like that
type Ast struct {
	structs
	globals
	functions map[string]Function
}

// something like this
type Parser struct {
	lexer *Lexer
	errors []string
}

// something like this
func (p *Parser) parse() *Ast{
	ast := new(Ast)
	for t := p.lexer.next(); t.ttype != "EOF"; t = p.lexer.next() {
		switch (t.ttype) {
		case "IMPORT": p.parseImport() // obviously not now, but this will prolly have to
									   // remake a lexer and all that shit
		case "STRUCT": p.parseStruct()
		case "GLOBAL": p.parseGlobal()
		case "FUNC":
			k, v := p.parseFunction(t)
			ast.functions[k] = v
		case "NEWLINE":
			continue
		}
	}
	if len(p.errors) != 0 {
		fmt.Println("Parsing error(s):")
		for _, s := range(p.errors) { fmt.Println(s) }
		os.Exit(0)
	}
	return ast
}

// something like that
func (p *Parser) parseFunction() (string, *Function) {
	f := new(Function{t: t})
	name, f.parameters, f.returnType := p.parseFunctionHeader()
	var block Block
	for t := p.lexer.next(); t.ttype != "}"; t = p.lexer.next() { // loop through statements
		block.appendStatement(p.parseStatement(t))
	}
	f.block = block // just append even if error occurs, doesn't matter since ast will
	return name, f  // never be returned if there exists an error
}
// couldve probably seperated this into some smaller functions but fuck it
// i'm really bad at this. So many conditionals, just don't see a way around it
// christ this function is lengthy. But I think it works
func (p *Parser) parseFunctionHeader() string { // not returning just a string
	t := p.lx.next()
	if t.ttype != "ID" {
		p.errorTrashLine(t, "Expecting function identifier on line %v after func", t.line)
		return
	}
	funcName := t.value
	if t2 := p.lx.next().value != "(" {
		p.errorTrashLine(t2, "Expecting '(' after %v on line %v", t.value, t.line)
		return
	}
	for t = p.lx.next() {
		if t.ttype != "TYPE" {
			p.errorTrashLine(t, "Expecting var declaration in function header, line %v", t.line)
			return
		}
		t2 := p.lx.next()
		if t2.ttype != "ID" {
			p.errorTrashLine(t2, "Expecting variable identifier after %v on line %v", t.value, t.line)
			return
		}
		// got a valid parameter, do something here
		t = p.lx.next()
		if t.value == ")" { break }
		if t.value == "," { continue }
		p.errorTrashLine(t, "Expexting ')' or ',' after %v on line %v", t2.value, t2.line)
		return // some bull shit
	}
	for t = p.lx.next(); t.value != "{"; t = p.lx.next() {
		if t.ttype != "TYPE" {
			p.errorTrashLine(t, "Return type examples: " +
					"') int, string {', ') bool {', ') {', etc. line %v", t.line)
			return
		}
		// add it to the return types
		if p.lx.peek().value == "," { p.lx.next() }
	}
	return funcName, /*parameters*/, /*return type(s), if any*/
}

//user by parse function
func (p *Parser) parseStatment(t token) Statement {
	switch (t.ttype) {
	case "TYPE": // for structs, might have to lex and parse and analze before funcs
		return p.parseDeclaration(t)
	case "KEYWORD": // if, for, etc
		return p.parseKeyword(t)
	case "ID":
		return p.parseReassignment(t)
	case "CALL":
		return p.parseFunctionCall(t)
	default:
		p.errorTrashLine(t, )
		return Statement{}
	}

}
// will be used by parseStatement (which is used by parse function) and parseGlobals
func (p *Parser) parseDeclaration(t) Declaration {
	d := &Declaration{t: t, ttype: t.value}
	if t = p.lx.next(); t.ttype != "ID" {
		p.errorTrashLine(t, "Expecting variable name after %v on line %v", d.ttype, t.line)
		return d,
	}
	d.id = t.value
	if t = p.lx.next(); t.value != "=" {
		p.errorTrashLine(t, "Expecting '=' after %v on line %v", d.id, t.line)
		return d
	}
	d.value = p.parseExpression()
	return d
}














func (p *Parser) errorTrashLine(t token, format string, args ...interface{}) {
	for t.ttype != "NEWLINE" {
		t = p.lx.next()
	}
	p.errors = append(p.errors, fmt.Sprintf(format, args...))
}

