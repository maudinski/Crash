package main
// TODO trashExpression really fucks up the rest of the errors
import (
    "strconv"
    "fmt"
    "os"
)

// http://eli.thegreenplace.net/2010/01/02/top-down-operator-precedence-parsing
// an article on pratt parsing, in case you forget what's going on
// passed 0 initially
func (p *Parser) parseExpression(rbp int) Expression {
    t := p.lx.next()
    f := p.nudFunctions[t.ttype]
    if f == nil {
        p.errorTrashExpression(t, "Invalid1 expression")
        return Id{}
    } // somehow need to account for EOF's and shit, logic through where to error check
    exp := f(p, t)
    for rbp < p.getBp(p.lx.peek().value) {
        t = p.lx.next()
        f2 := p.ledFunctions[t.value]
        if f2 == nil {
            p.errorTrashExpression(t, "Invalid 2Expression")
            return Id{}
        }
        exp = f2(p, t, exp) // t being passed solely for line numbers in lexical analysis
    }
    return exp
}

func (p *Parser) setPrattMaps() {
    p.nudFunctions = map[string]func(*Parser, token)Expression{"FLOAT_LITERAL": floatLiteral,
        "STRING_LITERAL": stringLiteral, "INT_LITERAL": intLiteral, "ID": id,
        "CALL": call, "BOOL_LITERAL": boolLiteral, "(": leftParen, "OPERATOR": infixOp}
    p.ledFunctions = map[string]func(*Parser, token, Expression)Expression{"+": add,
        "*": multiply, "/": divide, "-": subtract, "^": power, "==": equalequal,
        ">": greater, "<": less, ">=": greaterEqual, "<=": lessEqual, "!=": notEqual,
        "%": modulo}
    // changing these could fuck ip infixOp function, check
    p.bp = map[string]int{")": 0, "==": 10, ">": 10, "<": 10, ">=": 10, "<=": 10,
        "+": 20, "-": 20, "*": 30, "/": 30, "%": 30, "^": 40, "!=": 10, "&&": 5, "||": 5}
}

/*********nud functions************/
func floatLiteral(p *Parser, t token) Expression {
    f, _ := strconv.ParseFloat(t.value, 64) // i think 64 makes it a float64
    return Float{t, f}
}

func stringLiteral(p *Parser, t token) Expression {
    return String{t, t.value}
}

func intLiteral(p *Parser, t token) Expression {
    i, _ := strconv.Atoi(t.value)
    return Int{t, i}
}

func boolLiteral(p *Parser, t token) Expression {
    b, _ := strconv.ParseBool(t.value)
    return Bool{t, b} // TODO put this functionality into the lexer as a key word
}

func id(p *Parser, t token) Expression {
    return Id{t, t.value}
}

// probably
func call(p *Parser, t token) Expression {
    p.parsingFuncCallExpression++
    exp := p.parseFunctionCall(t)
    p.parsingFuncCallExpression--
    return exp
}

// this basically just starts parseExpression over, right after the left paren. Right
// parens bp is 0 so parseExpression will stop at that. Returns everything inside the
// parentheses
func leftParen(p *Parser, t token) Expression {
    // need to keep track that parenthesis are closed
    exp := p.parseExpression(0)
    if p.lx.peek().value != ")" {
        // there is mismatched parens
    } else {
        p.lx.next() // get rid of the right paren so it doesn't fuck shit up
    }
    return exp
}

// a switch for infix operators
// 50's are hardcoded to be higher than all other presedences for all operators. If you
// change the numbers in p.bp, then you might have to change these too
func infixOp(p *Parser, t token) Expression {
    switch (t.value) {
    case "-":
        return Negative{t: t, exp: p.parseExpression(50)}
    case "!":
        return Not{t: t, exp: p.parseExpression(50)}
    default:
        p.errorTrashExpression(t, "Invalid expression")
        return Id{}
    }
}

/*********************************/
/*************lef functioins*****************/
// all these functions and use of a map could seemingly be decomplicated with one function
// and a switch statement, but this is way is more efficient
func multiply(p *Parser, t token, e Expression) Expression {
    return Multiply{t: t, left: e, right: p.parseExpression(p.getBp(t.value))}
}

func add(p *Parser, t token, e Expression) Expression {
    return Add{t: t, left: e, right: p.parseExpression(p.getBp(t.value))}
}

func subtract(p *Parser, t token, e Expression) Expression {
    return Subtract{t: t, left: e, right: p.parseExpression(p.getBp(t.value))}
}

func divide(p *Parser, t token, e Expression) Expression {
    return Divide{t: t, left: e, right: p.parseExpression(p.getBp(t.value))}
}

func modulo(p *Parser, t token, e Expression) Expression {
    return Modulo{t: t, left: e, right: p.parseExpression(p.getBp(t.value))}
}

func power(p *Parser, t token, e Expression) Expression {
    return Power{t: t, base: e, exponent: p.parseExpression(p.getBp(t.value))}
}

func equalequal(p *Parser, t token, e Expression) Expression {
    return EqualEqual{t: t, left: e, right: p.parseExpression(p.getBp(t.value))}
}

func greater(p *Parser, t token, e Expression) Expression {
    return Greater{t: t, left: e, right: p.parseExpression(p.getBp(t.value))}
}

func less(p *Parser, t token, e Expression) Expression {
    return Less{t: t, left: e, right: p.parseExpression(p.getBp(t.value))}
}

func greaterEqual(p *Parser, t token, e Expression) Expression {
    return GreaterEq{t: t, left: e, right: p.parseExpression(p.getBp(t.value))}
}

func lessEqual(p *Parser, t token, e Expression) Expression {
    return LessEq{t: t, left: e, right: p.parseExpression(p.getBp(t.value))}
}

func notEqual(p *Parser, t token, e Expression) Expression {
    return NotEqual{t: t, left: e, right: p.parseExpression(p.getBp(t.value))}
}

func and(p *Parser, t token, e Expression) Expression {
    return And{t: t, left: e, right: p.parseExpression(p.getBp(t.value))}
}

func or(p *Parser, t token, e Expression) Expression {
    return Or{t: t, left: e, right: p.parseExpression(p.getBp(t.value))}
}

func (p *Parser) getBp(op string) int {
    bp, ok := p.bp[op]
    if !ok {
        return -1
    }
    return bp
}

// somethings fucked up with the way this works, so just gonna exit
// at the end of this and print the error. Bad but will atleast not give walk ass errors
// afterwords
// I think the fuck up has to do with the putBacks() and maybe with the way the lexers
// queue works, not sure though
func (p *Parser) errorTrashExpression(t token, msg string, args ...interface{}) {
    fullMsg := "Line " + toString(t.line) + ": " + msg
    for ; isExpressionPart(t); t = p.lx.next() {}
    if t.value != "," && p.parsingFuncCallExpression != 0 { // so parseFunctionCall doesnt
        p.lx.putBack(token{")", ")", t.line}) // fuck up
        fullMsg += " in function call"
    }
    p.lx.putBack(t)
    fmt.Printf(fullMsg + "\n", args...) // not supposed to be here
    os.Exit(0) // not supposed to be here
    //p.errors = append(p.errors, fmt.Sprintf(fullMsg, args...))
}

func isExpressionPart(t token) bool {
    switch (t.ttype) {
    case "INT_LITERAL": fallthrough
    case "FLOAT_LITERAL": fallthrough
    case "BOOL_LITERAL": fallthrough
    case "STRING_LITERAL": fallthrough
    case ")": fallthrough // gives some wrong errors if p.parserExp was called in parCall
    case "(": fallthrough
    case "OPERATOR": return true
    default: return false
    }
}


