package main

import (
	"fmt"
	"strconv"
)

// fucking AA if it's in the interface, then you don't have to do a type assertion :D
type Node interface {
	String() string
}

type Expression interface {
	Node
	//analyzeE(*SemAn) // used and defined in semanticAnalyzer.go
	isExpression() // dummy functions to force compiler to differentiate between E and S
}

type Statement interface {
	Node
	analyze(*SemAn) // used and defined in semanticAnalyzer.go
	isStatement()   // dummy functions to force compiler to differentiate between E and S
}

type Block []Statement

func newBlock() Block {
	b := make([]Statement, 0)
	return b
}
func (b Block) String() string {
	str := "Block: "
	str += "Statment amt = " + strconv.Itoa(len(b)) + "\n"
	for _, s := range b {
		str += "\t" + s.String() + "\n"
	}
	return str
}

/*******/
type If struct {
	t          token
	exp        Expression
	trueBlock  Block
	isElse     bool
	falseBlock Block
}

func (i If) isStatement() {}
func (i If) String() string {
	return fmt.Sprintf("\n\tIf:Expr: %v \n\tTrue block: %v \n\tFalse block: %v\n",
		i.exp.String(),
		i.trueBlock.String(),
		i.falseBlock.String())
}

/*******/
type Declaration struct {
	t     token
	ttype string
	id    Id
	value Expression // maybe not idk
}

func (d Declaration) isStatement() {}
func (d Declaration) String() string {
	return fmt.Sprintf("Declaration: %v %v = %v", d.ttype, d.id, d.value)
	//idk if this will work since Expression is an interface
}

type Function struct {
	t          token
	params     []Parameter
	returnType string
	name       string
	block      Block
}

func (f Function) String() string {
	str := fmt.Sprintf("Function name: %v | Parameters: ", f.name)
	for i, param := range f.params {
		str += strconv.Itoa(1+i) + ": " + param.String()
	}
	str += " | Return type: " + f.returnType
	return str + "\n" + f.block.String()
}

type Parameter struct {
	t     token
	ttype string
	id    Id
}

func (p Parameter) String() string {
	return fmt.Sprintf("type: %v, Id: %v ", p.ttype, p.id)
}

/****************/
type Id struct {
	t     token
	value string
}

func (id Id) isExpression() {}
func (id Id) String() string {
	return id.value
}

/*************/
type Reassignment struct {
	t     token
	id    Id
	value Expression
}

func (r Reassignment) isStatement()     {}
func (r Reassignment) String() string { // not sure if this will work since expression
	return fmt.Sprintf("Reassignment: Id %v, value %v", r.id, r.value) // is interface
}

/********************
type FakeExpression struct {
	t     token
	value string
}

func (fe FakeExpression) isExpression() {}
func (fe FakeExpression) String() string {
	return fe.value
}

/************/
type While struct {
	t         token
	condition Expression
	block     Block
}

func (w While) isStatement() {}
func (w While) String() string {
	return "While: condition: " + w.condition.String() + "block: " + w.block.String()
}

/************/
type Return struct {
	t     token
	value Expression // eventaully []Expression for multiple return values
}

func (r Return) isStatement() {}
func (r Return) String() string {
	return "Return " + r.value.String()
}

/**************/
type Call struct {
	t      token
	id     Id
	params []Expression
	// only if used in expression, this value is set in parser in parseExpression.
	typeShouldBe string
}

func (c Call) isStatement()  {}
func (c Call) isExpression() {} // Semantic analyzer will have to check if a function
//call used in an expression returns a value (And appropriate one but thats not the point)
func (c Call) String() string {
	str := "Function call, name " + c.id.String() + ", params: "
	for _, e := range c.params {
		str += e.String()
	}
	return str
}

/**************************this is all for expressions******************************88**/
type Int struct {
	t     token
	value int
}

func (i Int) isExpression() {}
func (i Int) String() string {
	return fmt.Sprintf("%v", i.value)
}

/**this would make sense right***/
type Float struct {
	t     token
	value float64
}

func (f Float) isExpression() {}
func (f Float) String() string {
	return strconv.FormatFloat(f.value, 'E', -1, 64) //idfk
}

/******this should make sense ****/
type String struct {
	t     token
	value string
}

func (s String) isExpression() {}
func (s String) String() string {
	return s.value
}

/********/
type Bool struct {
	t     token
	value bool
}
func (b Bool) isExpression() {}
func (b Bool) String() string {
	return strconv.FormatBool(b.value)
}
/*******/
type Not struct {
	t token
	exp Expression
}
func (n Not) isExpression() {}
func (n Not) String() string {
	return fmt.Sprintf("!(%v)", n.exp.String())
}
/*******/
type Negative struct {
	t token
	exp Expression
}
func (n Negative) isExpression() {}
func (n Negative) String() string {
	return fmt.Sprintf("-(%v)", n.exp.String())
}
/******************ooperation srcuts******************/
// a fuckton of structs but makes the code part easier with interfaces and implemented
// methods for compilation, checking validity, etc
/**********/
type Multiply struct {
	t     token
	left  Expression
	right Expression
}

func (m Multiply) isExpression() {}
func (m Multiply) String() string {
	return fmt.Sprintf("(%v * %v)", m.left.String(), m.right.String())
}

/**********/
type Add struct {
	t     token
	left  Expression
	right Expression
}

func (a Add) isExpression() {}
func (a Add) String() string {
	return fmt.Sprintf("(%v + %v)", a.left.String(), a.right.String())
}

/**********/
type Subtract struct {
	t     token
	left  Expression
	right Expression
}

func (s Subtract) isExpression() {}
func (s Subtract) String() string {
	return fmt.Sprintf("(%v - %v)", s.left.String(), s.right.String())
}

/**********/
type Divide struct {
	t     token
	left  Expression
	right Expression
}

func (d Divide) isExpression() {}
func (d Divide) String() string {
	return fmt.Sprintf("(%v / %v)", d.left.String(), d.right.String())
}

/**********/
type Modulo struct {
	t     token
	left  Expression
	right Expression
}

func (m Modulo) isExpression() {}
func (m Modulo) String() string {
	return fmt.Sprintf("(%v %% %v)", m.left.String(), m.right.String())
}

/**********/
type EqualEqual struct {
	t     token
	left  Expression
	right Expression
}

func (ee EqualEqual) isExpression() {}
func (ee EqualEqual) String() string {
	return fmt.Sprintf("(%v == %v)", ee.left.String(), ee.right.String())
}

/**********/
type NotEqual struct {
	t     token
	left  Expression
	right Expression
}

func (ne NotEqual) isExpression() {}
func (ne NotEqual) String() string {
	return fmt.Sprintf("(%v != %v)", ne.left.String(), ne.right.String())
}

/**********/
type Less struct {
	t     token
	left  Expression
	right Expression
}

func (l Less) isExpression() {}
func (l Less) String() string {
	return fmt.Sprintf("(%v < %v)", l.left.String(), l.right.String())
}

/**********/
type Greater struct {
	t     token
	left  Expression
	right Expression
}

func (g Greater) isExpression() {}
func (g Greater) String() string {
	return fmt.Sprintf("(%v > %v)", g.left.String(), g.right.String())
}

/**********/
type LessEq struct {
	t     token
	left  Expression
	right Expression
}

func (le LessEq) isExpression() {}
func (le LessEq) String() string {
	return fmt.Sprintf("(%v <= %v)", le.left.String(), le.right.String())
}

/**********/
type GreaterEq struct {
	t     token
	left  Expression
	right Expression
}

func (ge GreaterEq) isExpression() {}
func (ge GreaterEq) String() string {
	return fmt.Sprintf("(%v >= %v)", ge.left.String(), ge.right.String())
}

/*********/
type Power struct {
	t        token
	base     Expression
	exponent Expression
}

func (p Power) isExpression() {}
func (p Power) String() string {
	return fmt.Sprintf("(%v ^ %v)", p.base.String(), p.exponent.String())
}
