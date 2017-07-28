// this is shittily done. get it working then redo
package main

import (
	"fmt"
	"os"
)

// something like that
type Ast struct {
	structs []string // as in fuck this for now
	globals []Declaration
	functions map[string]Function
}

// something like this
type Parser struct {
	lx *Lexer
	errors []string
}

// something like this
// brain
func (p *Parser) parse() *Ast{
	ast := newAst() // these are all keywords but that's cool
	for t := p.lx.next(); t.ttype != "EOF"; t = p.lx.next() {
		switch (t.ttype) {
		case "IMPORT": p.parseImport() //prolly need to make a lexer, etc
		case "STRUCT": p.parseStruct()
		case "GLOBAL":
			Ast.globals = append(ast.globals, p.parseGlobal())
		case "FUNC":
			k, v := p.parseFunction(t)
			ast.functions[k] = v
		case "NEWLINE":
			continue
		default:
			p.errorTrashLine(t, "Unkown statement on line %v", t.line)
		}
	}
	if len(p.errors) != 0 {
		fmt.Println("Parsing error(s):")
		for _, s := range(p.errors) { fmt.Println(s) }
		os.Exit(0)
	}
	return ast
}
/****************GLOBAL SHIT******************/
func (p *Parser) parseGlobal() Declaration {
	t := p.lx.next() // doing this here cause parseDeclaration() takes the ttype token
	if t.ttype != "TYPE" { // through poor design/when parsing statements it's already
		p.errorTrashLine(t, "Expecting delaration after 'global' on line &v", t.line)// read
		return Declaration{}
	}
	return p.parseDeclaration(t)
}
/**************SOME UNHANDLED SHIT RIGHT NOW***********************/
func (p *Parser) parseImport() {
	fmt.Println("(p.parseImport) not handling import right now. Exit")
	os.Exit(1)
}

func (p *Parser) parseStruct() {
	fmt.Println("(p.parseStruct) Not parsing structs right now. Exit")
	os.Exit(1)
}
/****************PARSE FUNCTION SHIT******************************/
// something like that
func (p *Parser) parseFunction() (string, *Function) {
	f := new(Function{t: t})
	name, f.parameters, f.returnType := p.parseFunctionHeader()
	f.block = p.parseBlock() // go ahead everything takes care of it's own error
	return name, f  // never be returned if there exists an error (is this true? why
				//did I comment this?)
}
// couldve probably seperated this into some smaller functions but fuck it
// this is really bad rewrite this fucker
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

/****within function for the most part***/
/***SOME PARSE BLOCKS/STATEMENTS SHIT******/
func (p *Parser) parseBlock() *Block {
	b := newBlock()
	for t := p.lx.next(); t.value != "}"; t = p.lx.next() {
		if t.ttype == "NEWLINE" { continue }
		b.appendStatement(p.parseStatement())//TODO
	}
	return b
}


//user by parse function. brain for statment parsing
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
		p.errorTrashLine(t, "Not a valid statement on line %v", t.line)
		return Statement{}
	}

}

func (p *Parser) parseReassignment(t token) Reassignment {
	r := Reassignment{id: Id{value: t.value}}
	if p.lx.next().value != "=" {
		p.errorTrashLine(t, "Expecting '=' after %v for reassignmnet on line %v", t.value, t.line)
		return r
	}
	r.value = p.parseExpression()
	return r
}

func (p *Parser) parseKeyword(t token) Statement {
	switch (t.value) {
	case "if":
		return p.parseIf()
	case "for":
		return p.parseFor()
	case "return":
		return p.parseReturn()
	default:
		p.errorTrashLine(t, "Invalid use of keyword '%v' on line %v", t.value, t.line)
		return Statement{}
	}
}

// maybe not
func (p *Parser) parseFunctionCall(t token) Call{
	c := Call{id: Id{t.value}}
	c.params = make([]Expression, 0) // or however you do this
	//NOTE it's not possible for the next token to not be a (, because the lexer will
	//only provide a CALL type if it sees a (. So just consume it
	p.lx.next()
	if p.lx.peek().value == ")" { return c }
	err := false
	for t = p.lx.next() {
		c.params = append(c.params, p.parseExpression())
		t = p.lx.next()
		if t.value == ")" { break
		} else if t.value == "," { continue
		} else { err = true; break }
	}
	if err {
		p.errorTrashLine(t, "Expecting ')' or ',' in function call, got '%v'. Line %v", t.value, t.line)
	}
	return c
}

// will be used by parseStatement (which is used by parse function) and parseGlobals
func (p *Parser) parseDeclaration(t) Declaration {
	d := Declaration{t: t, ttype: t.value}
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

// I'm probably going to read this guys pratt parser shit, then just end up doing
// postfix expressions anyways, and do the lexical analysis part of expressions during
// this process. Should even be easier to generate code at runtime (maybe not actually)
// actually cause of function calls, seman will be harder during this. Maybe suck up a nut
// NOTE this will need to validate that it's starting an expression, nothing has looked
// at these tokens so far
// should also leave EVERYTHING after words intact, so if it reads in a comma, puts it
// the fuck back
func (p *Parser) parseExpression() /*notSure*/ {

}

func (p *Parser) errorTrashLine(t token, format string, args ...interface{}) {
	for t.ttype != "NEWLINE" && t.ttype != "EOF" {
		t = p.lx.next()
	}
	p.errors = append(p.errors, fmt.Sprintf(format, args...))
}

