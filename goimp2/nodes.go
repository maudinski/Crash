package main

import (
	"fmt"
	"strconv"
)
// fucking AA if it's in the interface, then you don't have to do a type assertion :D
type Node interface {
	String() string
}

type ExpressionPart interface {
	Node
	//analyzeE(*SemAn) // used and defined in semanticAnalyzer.go
	isExpressionPart() // dummy functions to force compiler to differentiate between E and S
}

type Statement interface {
	Node
	analyze(*SemAn) // used and defined in semanticAnalyzer.go
	isStatement() // dummy functions to force compiler to differentiate between E and S
}

// BUG BUG TODO this not being a pointer may be the source of a bug
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

/********************/
type FakeExpression struct {
	t     token
	value string
}

func (fe FakeExpression) isExpressionPart() {}
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
func (c Call) isExpressionPart() {} // Semantic analyzer will have to check if a function
//call used in an expression returns a value (And appropriate one but thats not the point)
func (c Call) String() string {
	str := "Function call, name " + c.id.String() + ", params: "
	for _, e := range c.params {
		str += e.String()
	}
	return str
}

type Postfix struct {
	t token
	exp []ExpressionPart
	overallType string
}
func (p Postfix) isExpressionPart() {}
func (p Postfix) String() string {
	str := ""
	for _, exp := range(p.exp) {
		str += exp.String()
	}
	return "PostfixExpression: " + str
}

type Operator struct {
	t token
	value string
}
func (o Operator) isExpressionPart() {}
func (o Operator) String() string {
	return o.value
}
/************************************************************************************8**
type InfixExpression struct {
	t token
	left Expression
	op string
	right Expression
}
func (ie InfixExpression) isExpression() {}
func (ie InfixExpression) String() string {
	return fmt.Sprintf("InfixExpressions:(%v %v %v)",
			ie.left.String(), ie.operator.String(), ie.right.String())
	//idk if this will work since expression is an interface
}
/*****/
type Int struct {
	t token
	value int
}
func (ie IntExpression) isExpressionPart() {}
func (ie IntExpression) String() string {
	return fmt.Sprintf("%v", ie.value)
}
/**this would make sense right***/
type Float struct{
	t token
	value float64
}
func (fe FloatExpression) isExpressionPart() {}
func (fe FloatExpression) String() string{
	return strconv.FormatFloat(fe.value, 'E', -1, 64) //idfk
}
/******this should make sense ****/
type String struct {
	t token
	value string
}
func (se StringExpressiong) isExpressionPart() {}
func (se StringExpression) String() string{
	return se.value
}
