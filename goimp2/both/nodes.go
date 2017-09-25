package main

// NOTE these are just intial defining of the nodes and their helper methods (like
// String, getToken, etc). They have more methods each. Semantic analyzer is just
// a few SemAn methods and the rest are someNode.analyze methods. code generator
// will probably be the same
import (
	"fmt"
	"strconv"
)

// ___ing AA if it's in the interface, then you don't have to do a type ___ertion :D
type Node interface {
	String() string
}

type Expression interface {
	Node
	analyzeE(*SemAn) string // used and defined in semanticAnalyzer.go
	isExpression()          // dummy functions to force compiler to differentiate between E and S
	getT() token            // is expression kinda pointless now but ohwell. Nice and readable
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
	falseBlock Block // golangs parser has blocks as statements, so that this falseBlock
}	// (namely, 'else' block) just takes a statement, which means this else block can
	// be another if statement, allowing if else chains
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
	str := "Function call, name " + c.id.String() + ", params:"
	for i, e := range c.params {
		str += " |" + toString(i+1) + ":" + e.String()
	}
	return str
}
func (a Call) getT() token { return a.t } // a for m___ copy and paste
/**************************this is all for expressions******************************88**/
type Id struct {
	t     token
	value string
}

func (id Id) isExpression() {}
func (id Id) String() string {
	return id.value
}
func (a Id) getT() token { return a.t } // a for m___ copy and paste

/*************/
type Int struct {
	t     token
	value int
}

func (i Int) isExpression() {}
func (i Int) String() string {
	return fmt.Sprintf("%v", i.value)
}
func (a Int) getT() token { return a.t } // a for m___ copy and paste

/**this would make sense right***/
type Float struct {
	t     token
	value float64
}

func (f Float) isExpression() {}
func (f Float) String() string {
	return strconv.FormatFloat(f.value, 'E', -1, 64) //idfk
}
func (a Float) getT() token { return a.t } // a for m___ copy and paste

/******this should make sense ****/
type String struct {
	t     token
	value string
}

func (s String) isExpression() {}
func (s String) String() string {
	return s.value
}
func (a String) getT() token { return a.t } // a for m___ copy and paste

/********/
type Bool struct {
	t     token
	value bool
}

func (b Bool) isExpression() {}
func (b Bool) String() string {
	return strconv.FormatBool(b.value)
}
func (a Bool) getT() token { return a.t } // a for m___ copy and paste

/*******/
type Not struct {
	t   token
	exp Expression
}

func (n Not) isExpression() {}
func (n Not) String() string {
	return fmt.Sprintf("!(%v)", n.exp.String())
}
func (a Not) getT() token { return a.t } // a for m___ copy and paste

/*******/
type Negative struct {
	t   token
	exp Expression
}

func (n Negative) isExpression() {}
func (n Negative) String() string {
	return fmt.Sprintf("-(%v)", n.exp.String())
}
func (a Negative) getT() token { return a.t } // a for m___ copy and paste

/******************ooperation srcuts******************/
// a ___ton of structs but makes the code part easier with interfaces and implemented
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
func (a Multiply) getT() token { return a.t } // a for m___ copy and paste

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
func (a Add) getT() token { return a.t } // a for m___ copy and paste

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
func (a Subtract) getT() token { return a.t } // a for m___ copy and paste

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
func (a Divide) getT() token { return a.t } // a for m___ copy and paste

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
func (a Modulo) getT() token { return a.t } // a for m___ copy and paste

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
func (a EqualEqual) getT() token { return a.t } // a for m___ copy and paste

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
func (a NotEqual) getT() token { return a.t } // a for m___ copy and paste

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
func (a Less) getT() token { return a.t } // a for m___ copy and paste

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
func (a Greater) getT() token { return a.t } // a for m___ copy and paste

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
func (a LessEq) getT() token { return a.t } // a for m___ copy and paste

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
func (a GreaterEq) getT() token { return a.t } // a for m___ copy and paste

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
func (a Power) getT() token { return a.t } // a for m___ copy and paste

/**********/
type And struct {
	t     token
	left  Expression
	right Expression
}

func (a And) isExpression() {}
func (a And) String() string {
	return fmt.Sprintf("(%v >= %v)", a.left.String(), a.right.String())
}
func (a And) getT() token { return a.t } // a for m___ copy and paste

/**********/
type Or struct {
	t     token
	left  Expression
	right Expression
}

func (o Or) isExpression() {}
func (o Or) String() string {
	return fmt.Sprintf("(%v >= %v)", o.left.String(), o.right.String())
}
func (a Or) getT() token { return a.t } // a for m___ copy and paste
