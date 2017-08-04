// this is hacked get it working then redo
package main

// TODO error message to format line12: errMsg
// TODO got and comment what token each function is taking as a parameter and some
// example of what they are parsing, I'll thank myself later I'm sure
import (
	"fmt"
	"os"
)

// something like this
type Parser struct {
	lx     *Lexer
	errors []string
	// nud stands for null denotion, which in contrast to the beneath one, mean "these
	// are the functions im looking for when i have no unclaimed expressions to tie this
	// to". Refers to single expressions (like literals and ID's, function calls, etc)
	// and also infix operators (like ! and negative sign). Those are the keys
	//  set with p.setPrattMaps(), in expressionParser.go
	nudFunctions map[string]func(*Parser, token) Expression
	// led stands for left denotion, whcih i guess means that "these are the functions
	// im looking for when i have some left expression that is unclaimed". These usually
	// refer to operators. Those are the key
	// https://people.csail.mit.edu/jaffer/slib/Nud-and-Led-Definition.html
	// link to some article. set with p.setPrattMaps()
	ledFunctions map[string]func(*Parser, token, Expression) Expression
	// stands for binding power, which is literally just presedence power, but pratt
	// describes it as the power an operator has of binding an expression to it. ie:
	// 2 + 3 * 5, the 3 gets bound to the *, not the +, cause of presedence. Turns to
	// 2 + (3 * 5)  set with p.setPrattMaps()
	bp map[string]int
	// pretty nasty work around, but this is set and unset in call() nud function for
	// pratt parser. It's used in errorTrashExpression() to put back a ) if the error
	// ocurred in an expression in a function call, so that parseFunctionCall() can give
	// accurate errors
	parsingFuncCallExpression int
}

func newParser(lx *Lexer) *Parser {
	p := new(Parser)
	p.lx = lx
	p.parsingFuncCallExpression = 0
	p.errors = make([]string, 0)
	return p
}

// something like this
// brain
func (p *Parser) parse() *Ast {
	ast := newAst() // these are all keywords but that's cool
	for t := p.lx.next(); t.ttype != "EOF"; t = p.lx.next() {
		switch t.value {
		case "import":
			p.parseImport() //prolly need to make a lexer, etc
		case "struct":
			p.parseStruct()
		case "global":
			ast.globals = append(ast.globals, p.parseGlobal())
		case "func":
			k, v := p.parseFunction(t)
			if _, exists := ast.functions[k]; exists {
				p.errorTrashLine(t, "Function '%v' already defined", k)
			}
			ast.functions[k] = v
		case "\\n":
			continue
		default:
			p.errorTrashLine(t, "Token/statment outside function. '%v'...", t.value)
		}
	}
	if len(p.errors) != 0 {
		fmt.Println("Parsing error(s):")
		for _, s := range p.errors {
			fmt.Println(s)
		}
		os.Exit(0)
	}
	return ast
}

/****************GLOBAL SHIT******************/
func (p *Parser) parseGlobal() Declaration {
	t := p.lx.next()       // doing this here cause parseDeclaration() takes the ttype token
	if t.ttype != "TYPE" { // through poor design/when parsing statements it's already
		p.errorTrashLine(t, "Expecting delaration after 'global'") // read
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
	f.name, f.params, f.returnType = p.parseFunctionHeader()
	f.block = p.parseBlock() // go ahead everything takes care of it's own error
	return f.name, f         // never be returned if there exists an error (is this true? why
	//did I comment this?)
}

// couldve probably seperated this into some smaller functions but fuck it
// this is really bad rewrite this fucker
// christ this function is lengthy. But I think it works
// ma fucking cheted. Works tho
// TODO rewrite
func (p *Parser) parseFunctionHeader() (string, []Parameter, string) {
	t := p.lx.next()
	params := make([]Parameter, 0)
	if t.ttype != "ID" {
		p.errorTrashLine(t, "Expecting function identifier after func")
		return "", params, ""
	}
	funcName := t.value
	if t2 := p.lx.next(); t2.value != "(" {
		p.errorTrashLine(t2, "Expecting '(' after %v", t.value)
		return funcName, params, ""
	}
	if t = p.lx.next(); t.value != ")" {
		for ; true; t = p.lx.next() { // exiting handled in loop. I prefer while loops
			if t.ttype != "TYPE" {
				p.errorTrashLine(t, "Expecting var declaration in function header")
				return funcName, params, ""
			}
			t2 := p.lx.next()
			if t2.ttype != "ID" {
				p.errorTrashLine(t2, "Expecting variable identifier after %v", t.value)
				return funcName, params, ""
			}
			params = append(params, Parameter{t2, t.value, Id{t2, t2.value}})
			t = p.lx.next()
			if t.value == ")" {
				break
			}
			if t.value == "," {
				continue
			}
			p.errorTrashLine(t, "Expexting ')' or ',' after %v", t2.value)
			return funcName, params, ""
		}
	} /* should work but only doing one return type for now
	for t = p.lx.next(); t.value != "{"; t = p.lx.next() {
		if t.ttype != "TYPE" {
			p.errorTrashLine(t, "Return type examples: " +
					"') int, string {', ') bool {', ') {', etc. line %v", t.line)
			return
		}
		// add it to the return types
		if p.lx.peek().value == "," { p.lx.next() }
	}*/
	/***this is for one return type***/
	t = p.lx.next()
	returnType := ""
	if t.ttype == "TYPE" {
		returnType = t.value
	} else if t.value == "{" {
		p.lx.putBack(t)
	} else {
		p.errorTrashLine(t, "Expecting {, got %v", t.value)
	}
	return funcName, params, returnType
}

/****within function for the most part***/
/***SOME PARSE BLOCKS/STATEMENTS SHIT******/
func (p *Parser) parseBlock() Block {
	b := newBlock()
	eof := false
	var t1 token
	for t1 = p.lx.next(); t1.value != "{"; {
		p.errorTrashLine(t1, "Expecting {, got %v", t1.value)
	}
	for t := p.lx.next(); t.value != "}"; t = p.lx.next() {
		if t.ttype == "NEWLINE" {
			continue
		}
		if t.ttype == "EOF" {
			eof = true
			break
		}
		b = append(b, p.parseStatement(t))
	}
	if eof {
		p.errorTrashLine(t1, "Block never closed. Need }")
	}
	return b
}

//user by parse function. brain for statment parsing
func (p *Parser) parseStatement(t token) Statement {
	switch t.ttype {
	case "TYPE": // for structs, might have to lex and parse and analze before funcs
		return p.parseDeclaration(t)
	case "KEYWORD": // if, for, etc
		return p.parseKeyword(t)
	case "ID":
		return p.parseReassignment(t)
	case "CALL":
		return p.parseFunctionCall(t)
	default:
		p.errorTrashLine(t, "Not a valid statement")
		return Declaration{} // just need to return somestatment, arbitrary
	} // parser will exit if any errors exist

}

func (p *Parser) parseReassignment(t token) Reassignment {
	r := Reassignment{t: t, id: Id{value: t.value}}
	if p.lx.next().value != "=" {
		p.errorTrashLine(t, "Expecting '=' after %v for reassignmnet", t.value)
		return r
	}
	r.value = p.parseExpression(0)
	return r
}

// maybe not
func (p *Parser) parseFunctionCall(t token) Call {
	c := Call{t: t, id: Id{t, t.value}}
	c.params = make([]Expression, 0) // or however you do this
	//it's not possible for the next token to not be a (, because the lexer will
	//only provide a CALL type if it sees a (. So just consume it
	p.lx.next()
	t = p.lx.next()
	if t.value == ")" {
		return c
	}
	p.lx.putBack(t)
	err := false
	c.params = append(c.params, p.parseExpression(0))
	for t = p.lx.next(); true; t = p.lx.next() { // break logic is handled in loop
		if t.value == ")" {
			break
		} else if t.value == "," {
			c.params = append(c.params, p.parseExpression(0))
		} else {
			err = true
			break
		}
	}
	if err {
		p.errorTrashLine(t, "Expecting ')' or ',' in function call, got '%v'", t.value)
	}
	return c
}

// will be used by parseStatement (which is used by parse function) and parseGlobals
func (p *Parser) parseDeclaration(t token) Declaration {
	d := Declaration{t: t, ttype: t.value}
	if t = p.lx.next(); t.ttype != "ID" {
		p.errorTrashLine(t, "Expecting variable name after %v", d.ttype)
		return d
	}
	d.id = Id{t, t.value}
	if t = p.lx.next(); t.value != "=" {
		p.errorTrashLine(t, "Expecting '=' after %v", d.id)
		return d
	}
	d.value = p.parseExpression(0)
	return d
}

/*************KEYWORD PARSING****************/
func (p *Parser) parseKeyword(t token) Statement {
	switch t.value {
	case "if":
		return p.parseIf(t)
	case "while":
		return p.parseWhile(t)
	case "return":
		return p.parseReturn(t)
	default:
		p.errorTrashLine(t, "Invalid use of keyword '%v'", t.value)
		return Declaration{} // arbitrary, just need to return something
	} // parser will exit if any errors exist
}

func (p *Parser) parseIf(t token) If {
	i := If{t: t, exp: p.parseExpression(0), trueBlock: p.parseBlock()}
	if t = p.lx.next(); t.value == "else" {
		i.isElse = true
		i.falseBlock = p.parseBlock()
	} else {
		i.isElse = false
	}
	return i
}

func (p *Parser) parseWhile(t token) While {
	return While{t, p.parseExpression(0), p.parseBlock()}
}

//will only work for single return types
func (p *Parser) parseReturn(t token) Return {
	return Return{t, p.parseExpression(0)}
}

// dirty-ish but oh well
func (p *Parser) errorTrashLine(t token, format string, args ...interface{}) {
	fullMsg := "Line " + toString(t.line) + ": " + format
	leftBrace := false
	for t.ttype != "NEWLINE" && t.ttype != "EOF" {
		if t.ttype == "{" {
			leftBrace = true
		}
		t = p.lx.next()
	}
	if leftBrace && t.ttype != "EOF" {
		p.lx.putBack(token{"{", "{", t.line})
	}
	fmt.Printf(fullMsg+"\n", args...) // added this is cause queueing up  relvant errors
	os.Exit(0)                        // sucks
	p.errors = append(p.errors, fmt.Sprintf(fullMsg, args...))
}
