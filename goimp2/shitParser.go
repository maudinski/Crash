// this is shittily done. get it working then redo
package main

import (
	"fmt"
	"os"
)

// something like this
type Parser struct {
	lx *Lexer
	errors []string
}

func newParser(lx *Lexer) *Parser {
	p := new(Parser)
	p.lx = lx
	p.errors = make([]string, 0)
	return p
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
			ast.globals = append(ast.globals, p.parseGlobal())
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
func (p *Parser) parseFunction(t token) (string, *Function) {
	f := new(Function)
	f.t = t
	var name string
	name, f.params, f.returnType = p.parseFunctionHeader()
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
	return funcName,
}

/****within function for the most part***/
/***SOME PARSE BLOCKS/STATEMENTS SHIT******/
func (p *Parser) parseBlock() Block {
	b := newBlock()
	eof := false
	var t1 token
	if t1 = p.lx.next(); t1.value != "{" {
		p.errorTrashLine(t1, "Expecting { on line %v, got %v", t1.line, t1.value)
	}
	for t := p.lx.next(); t.value != "}"; t = p.lx.next() {
		if t.ttype == "NEWLINE" { continue }
		if t.ttype == "EOF" { eof = true; break }
		b.appendStatement(p.parseStatement(t))
	}
	if eof {
		p.errorTrashLine(t1, "Block never closed on line %v. Need }", t1.line)
	}
	return b
}

//user by parse function. brain for statment parsing
func (p *Parser) parseStatement(t token) Statement {
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
		return Declaration{} // just need to return somestatment, arbitrary
	}	// parser will exit if any errors exist

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

// maybe not
func (p *Parser) parseFunctionCall(t token) Call{
	c := Call{t: t, id: Id{t, t.value}}
	c.params = make([]Expression, 0) // or however you do this
	//NOTE it's not possible for the next token to not be a (, because the lexer will
	//only provide a CALL type if it sees a (. So just consume it
	p.lx.next()
	if p.lx.peek().value == ")" { return c }
	err := false
	for t = p.lx.next(); true; t = p.lx.next() { // break logic is handled in loop
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
func (p *Parser) parseDeclaration(t token) Declaration {
	d := Declaration{t: t, ttype: t.value}
	if t = p.lx.next(); t.ttype != "ID" {
		p.errorTrashLine(t, "Expecting variable name after %v on line %v", d.ttype, t.line)
		return d
	}
	d.id = Id{t, t.value}
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
// when written, could probably check that the next token is something valid? say, has
// to be a a comma, ), newline, or }
func (p *Parser) parseExpression() Expression {
	//TODO, this just sets up a fake expression. For debugging
	exp := ""
	for t := p.lx.next(); true; t = p.lx.next() { // logic handled in loop. dummy code anyways
		if t.value == ")" || t.value =="," || t.value == "\\n" { break }
		exp += t.value
	}
	return FakeExpression{value: exp}
}

/*************KEYWORD PARSING****************/
func (p *Parser) parseKeyword(t token) Statement {
	switch (t.value) {
	case "if":
		return p.parseIf(t)
	case "while":
		return p.parseWhile(t)
	case "return":
		return p.parseReturn(t)
	default:
		p.errorTrashLine(t, "Invalid use of keyword '%v' on line %v", t.value, t.line)
		return Declaration{} // arbitrary, just need to return something
	}	// parser will exit if any errors exist
}
// TODO go through and add tokens to everything
func (p *Parser) parseIf(t token) If {
	i := If{t: t, exp: p.parseExpression(), trueBlock: p.parseBlock()}
	if t = p.lx.next(); t.value == "else" {
		i.isElse = true
		i.falseBlock = p.parseBlock()
	} else {
		i.isElse = false
	}
	return i
}

func (p *Parser) parseWhile(t token) While {
	return While{t, p.parseExpression(), p.parseBlock()}
}

//will only work for single return types
func (p *Parser) parseReturn(t token) Return {
	return Return{t, p.parseExpression()}
}

func (p *Parser) errorTrashLine(t token, format string, args ...interface{}) {
	for t.ttype != "NEWLINE" && t.ttype != "EOF" {
		t = p.lx.next()
	}
	p.errors = append(p.errors, fmt.Sprintf(format, args...))
}

