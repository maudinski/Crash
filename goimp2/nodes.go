package main

import(
	"fmt"
	"strconv"
)

type Node interface {
	String() string
}

type Expression interface {
	Node
	isExpression() // dummy functions to force compiler to differentiate between E and S
}

type Statement interface {
	Node
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
	str += "Statment amt = " + strconv.Itoa(len(b)) +"\n"
	for _, state := range(b) {
		switch s := state.(type) {
		default:
			str += "\t" + s.String() + "\n"
		} // a tour of go i suggesting that s will be the approriate variable...
	}	// in which case this should work in all cases. Jesus I fucking hope
	return str // it does :)
}
/*******/
type If struct {
	t token
	exp Expression
	trueBlock Block
	isElse bool
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
	t token
	ttype string
	id Id
	value Expression// maybe not idk
}
func (d Declaration) isStatement() {}
func (d Declaration) String() string {
	return fmt.Sprintf("Declaration: %v %v = %v", d.ttype, d.id, d.value)
	//idk if this will work since Expression is an interface
}
/*****
type InfixExpression struct {
	t token
	left Expression
	operator string
	right Expression
}
func (ie InfixExpression) isExpression() {}
func (ie InfixExpression) String() string {
	return fmt.Sprintf("InfixExpressions:(%v %v %v)", ie.left, ie.operator, ie.right)
	//idk if this will work since expression is an interface
}
/*****
type IntExpression struct {
	t token
	value int
}
func (ie IntExpression) isExpression() {}
func (ie IntExpression) String() string {
	return fmt.Sprintf("%v", ie.value)
}
/**this would make sense right***
type FloatExpression struct{
	t token
	value float64
}
func (fe FloatExpression) isExpression() {}
func (fe FloatExpression) String() string{
	return string(fe.value)
}
/******this should make sense ****
type StringExpression struct {
	t token
	value string
}
func (se StringExpressiong) isExpression() {}
func (se StringExpression) String() string{
	return se.value
}


/********/
//idk about this yet
type Function struct {
	t token
	params []Parameter
	returnType  string
	name string
	block Block
}
func (f Function) String() string {
	str := fmt.Sprintf("Function name: %v | Parameters: ", f.name)
	for i, param := range(f.params) {
		str += strconv.Itoa(1+i) + ": " + param.String()
	}
	str += " | Return type: " + f.returnType
	return str + "\n" + f.block.String()
}

type Parameter struct {
	t token
	ttype string
	id Id
}
func (p Parameter) String() string {
	return fmt.Sprintf("type: %v, Id: %v ", p.ttype, p.id)
}
/****************/
type Id struct {
	t token
	value string
}
func (id Id) isExpression() {}
func (id Id) String() string {
	return id.value
}
/*************/
type Reassignment struct {
	t token
	id Id
	value Expression
}
func (r Reassignment) isStatement() {}
func (r Reassignment) String() string { // not sure if this will work since expression
	return fmt.Sprintf("Reassignment: Id %v, value %v", r.id, r.value)// is interface
}

/**************/
type Call struct {
	t token
	id Id
	params []Expression
}
func (c Call) isStatement() {}
func (c Call) isExpression() {} // Semantic analyzer will have to check if a function
//call used in an expression returns a value (And appropriate one but thats not the point)
func (c Call) String() string {
	str := "Function call, name " + c.id.String() + ", params: "
	for _, e := range(c.params) {
		switch exp := e.(type) {
		default:
			str += exp.String()
		}
	}
	return str
}
/********************/
type FakeExpression struct {
	t token
	value string
}
func (fe FakeExpression) isExpression(){}
func (fe FakeExpression) String() string {
	return fe.value
}
/************/
type While struct {
	t token
	condition Expression
	block Block
}
func (w While) isStatement() {}
func (w While) String() string {
	return "While: condition: " + w.condition.String() + "block: " + w.block.String()
}
/************/
type Return struct {
	t token
	value Expression // eventaully []Expression for multiple return values
}
func (r Return) isStatement() {}
func (r Return) String() string {
	return "Return " + r.value.String()
}


































