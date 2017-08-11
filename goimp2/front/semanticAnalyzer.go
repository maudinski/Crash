// TODO need to create a symbol table that stores the symbols and types (types or size of
// memory needed) and also has a spot for the relative address it will be at. Not sure
// about relative address being in there. 9.3 in dragonbook for address shit
package main
// point of this is to check if variables exist in the scope, functions called where
// already defined/built in, expressions are valid, and type checking for expression
// and variable assignments/reassignments. Final string of error checking
// NOTE expression functions are excessively redundant but it's easier/faster in the end
// that way. Bth for semantic analysis and parsing is true
import (
	"fmt"
	"os"
)

type DeclaredFunc struct {
	name       string   // redundent cause it's being stored in a mao by it's name but ohwell
	params     []string // slice of the types
	returnType string   // will eventually be able to return multiple types
}

type SemAn struct {
	ast    *Ast
	es     *EnvStack
	errors []string
	// phase1. Holds info from function headers. Only for a more convenient format
	// than the []*Function in ast. Key is the name of the function
	declaredFuncs map[string]DeclaredFunc
	// holds the name of the current function being parsed
	// for analyzing return statments
	currentFunc string
}

func newSemAn(ast *Ast) *SemAn {
	sa := new(SemAn)
	sa.ast = ast
	sa.es = newEnvironmentStack()
	sa.errors = make([]string, 0)
	sa.declaredFuncs = make(map[string]DeclaredFunc, 0)
	return sa
}

// this doesn't actually modify the ast at all
// phase sounds a lot cooler than it is
func (sa *SemAn) analyze() {
	sa.phase1()
	sa.checkErrors() // will only be global variable expression errors
	sa.phase2()
	sa.checkErrors()
	sa.phase3() // does it's own exiting
	sa.checkErrors()
}

// pre shit
// 1) pulls globals from sa.ast.globals
// 2) pulls defined functions from ast.functions, they're parameter types in order,
//    and the return type(s) in more accessible manner
// phase really makes this sound way more epic than it is
// might be good
func (sa *SemAn) phase1() {
	sa.es.pushNewEnv()
	for _, dec := range sa.ast.globals { // add global env to stack, never popped
		dec.analyze(sa) // remeber that control switches to the statement, no confuse
	} // analyze already adds it to the current environment
	for n, f := range sa.ast.functions { // name, *Function
		params := make([]string, 0) // get a []string of the parameter types
		for _, p := range f.params {
			params = append(params, p.ttype) // slice of the types
		}
		sa.declaredFuncs[n] = DeclaredFunc{n, params, f.returnType}
	}

}

// calls analyze block on each function in the ast
func (sa *SemAn) phase2() {
	for name, f := range sa.ast.functions {
		sa.es.pushNewEnv()           // new env created before each block is entered
		for _, p := range f.params { // add all the parameters to the environement
			sa.es.add(p.id.value, p.ttype)
		}
		sa.currentFunc = name
		sa.analyzeBlock(f.block)
		sa.es.popEnv() // pop the environment
	}
}

// for now just checks that a main function exists. Can be for any more small checks that
// are needed
func (sa *Seman) phase3() {
	_, ok := sa.declaredFuncs["main"]
	if !ok {
		fmt.Println("No main function, compilation failed")
		os.Exit(0)
	}
}

/******this switches control over to the statement********/
func (sa *SemAn) analyzeBlock(block Block) {
	for _, s := range block {
		s.analyze(sa)
	}
}

/********recievers from here on out turn to the statement and get sa passed**********/
/***turns out to a little less code in parse block, and eaiser modification**********/

// something like that
func (d Declaration) analyze(sa *SemAn) {
	exists, _ := sa.es.checkTop(d.id.value)
	if exists {
		sa.error(d.t, "Var %v already exists in this scope", d.id.value)
		return
	}
	sa.es.add(d.id.value, d.ttype)
	sa.analyzeExpression(d.ttype, d.value)
}

// analyzeExpression verifies the expression and makes sure the resulting type
// is of passed in type
// so pass it the expression and pass it the return type it's supposed to be
func (r Return) analyze(sa *SemAn) {
	sa.analyzeExpression(sa.declaredFuncs[sa.currentFunc].returnType, r.value)
}

func (i If) analyze(sa *SemAn) {
	sa.analyzeExpression("bool", i.exp)
	sa.es.pushNewEnv()
	sa.analyzeBlock(i.trueBlock)
	sa.es.popEnv()
	// i might be stupid but I can't understand or remember why this is !i.isElse, since
	// the oposite would make more sense. But it works, so fuck it who cares
	if !i.isElse {
		sa.es.pushNewEnv()
		sa.analyzeBlock(i.falseBlock)
		sa.es.popEnv()
	}
}

func (w While) analyze(sa *SemAn) {
	sa.analyzeExpression("bool", w.condition)
	sa.es.pushNewEnv()
	sa.analyzeBlock(w.block)
	sa.es.popEnv()
}

// idk why i did this or where the trasnfer is but the types are lowercase
func (r Reassignment) analyze(sa *SemAn) {
	exists, ttype := sa.es.check(r.id.value)
	if !exists {
		sa.error(r.t, "Var %v not declared", r.id.value)
		return
	}
	if ttype == "string" {
		sa.error(r.t, "Cannot reassign string %v. Immutable", r.id.value)
		return
	}
	sa.analyzeExpression(ttype, r.value)
}

func (c Call) analyze(sa *SemAn) {
	declaredFunc, ok := sa.declaredFuncs[c.id.value]
	if !ok { // in here i guess check if it's a built in
		sa.error(c.t, "Function %v not defined", c.id.value)
		return
	}
	if len(declaredFunc.params) != len(c.params) {
		sa.error(c.t, "Wrong number of arguments to function call %v", c.id.value)
		return
	}
	for i, exp := range c.params {
		sa.analyzeExpression(declaredFunc.params[i], exp)
	}
}

/*********this goes back to sa reciever but switches it over to expression*********/
// should validate the expression and make sure it is of the same type passed
// something like that
// ... should be good
func (sa *SemAn) analyzeExpression(supposedType string, e Expression) {
	if ttype := e.analyzeE(sa); supposedType != ttype {
		sa.error(e.getT(), "Mismatched types. Want %v, have %v", supposedType, ttype)
	}
}

/*************reciever goes over to expression and takes sa as param*************/
/*************there could definitely be some less redundent ways of doing this but this
  NOTE       takes advantage of the speed/simplicity of interface methods*************/
// thinking these return the type of the expression
func (id Id) analyzeE(sa *SemAn) string {
	exists, ttype := sa.es.check(id.value)
	if !exists {
		sa.error(id.t, "Variable %v not declared", id.value)
	}
	return ttype
}

// basically be c.analyze(sa) but with a return type. just copied the code over
func (c Call) analyzeE(sa *SemAn) string {
	declaredFunc, exists := sa.declaredFuncs[c.id.value]
	if !exists {
		sa.error(c.t, "Function %v not declared", c.id.value)
		return ""
	} else if len(declaredFunc.params) != len(c.params) {
		sa.error(c.t, "Wrong number of arguments to function call %v", c.id.value)
		return declaredFunc.returnType
	}
	for i, exp := range c.params {
		sa.analyzeExpression(declaredFunc.params[i], exp)
	}
	return declaredFunc.returnType
}

// the literals will already for a fact be ok, since they are determined in the lexer
// and set accordingly in the pratt parser machine
func (i Int) analyzeE(sa *SemAn) string {
	return "int"
}

func (s String) analyzeE(sa *SemAn) string {
	return "string"
}

func (f Float) analyzeE(sa *SemAn) string {
	return "float"
}

func (b Bool) analyzeE(sa *SemAn) string {
	return "bool"
}

func (n Not) analyzeE(sa *SemAn) string {
	if ttype := n.exp.analyzeE(sa); ttype != "bool" {
		sa.error(n.t, "'!' operator needs boolean value", ttype)
	}
	return "bool"
}

func (n Negative) analyzeE(sa *SemAn) string {
	ttype := n.exp.analyzeE(sa)
	if ttype != "int" && ttype != "float" {
		sa.error(n.t, "'-' needs int or float value, got %v", ttype)
	}
	return ttype
}

func (m Multiply) analyzeE(sa *SemAn) string {
	left := m.left.analyzeE(sa)
	right := m.right.analyzeE(sa)
	if left != right {
		sa.error(m.t, "Mismatched types for '*'. %v * %v", left, right)
	} else if left != "int" && left != "float" { // doesnt matter left or right
		sa.error(m.t, "'*' needs int or float values, have %v and %v", left, right)
	}
	return left // doesnt matter left or right
}

// 1) doing all three comparisons instead of just left == "bool" because structs will
// be tossed intol lexers.addTypes as they are encountered so they will also be types
// also, all left comparisons there cause it doesnt matter, left == right at that point
func (a Add) analyzeE(sa *SemAn) string {
	left := a.left.analyzeE(sa)
	right := a.right.analyzeE(sa)
	if left != right {
		sa.error(a.t, "Mismatched types for '+'. %v + %v", left, right)
	} else if left != "int" && left != "float" && left != "string" {
		sa.error(a.t, "'+' needs int, string, or float values, have %v and %v", left, right)
	}
	return left // doesnt matter left or right
}

func (s Subtract) analyzeE(sa *SemAn) string {
	left := s.left.analyzeE(sa)
	right := s.right.analyzeE(sa)
	if left != right {
		sa.error(s.t, "Mismatched types for '-'. %v - %v", left, right)
	} else if left != "int" && left != "float" { // doesnt matter left or right
		sa.error(s.t, "'-' needs int or float values, have %v and %v", left, right)
	}
	return left // doesnt matter left or right
}

func (d Divide) analyzeE(sa *SemAn) string {
	left := d.left.analyzeE(sa)
	right := d.right.analyzeE(sa)
	if left != right {
		sa.error(d.t, "Mismatched types for '/'. %v / %v", left, right)
	} else if left != "int" && left != "float" { // doesnt matter left or right
		sa.error(d.t, "'/' needs int or float values, have %v and %v", left, right)
	}
	return left // doesnt matter left or right
}

func (m Modulo) analyzeE(sa *SemAn) string {
	left := m.left.analyzeE(sa)
	right := m.right.analyzeE(sa)
	if left != right {
		sa.error(m.t, "Mismatched types for '%'. %v %% %v", left, right)
	} else if left != "int" { // doesnt matter left or right
		sa.error(m.t, "'%' needs int values, have %v and %v", left, right)
	}
	return left // doesnt matter left or right
}

// NOTE starting to use b for all recievers cause copy and paste
func (b EqualEqual) analyzeE(sa *SemAn) string {
	left := b.left.analyzeE(sa)
	right := b.right.analyzeE(sa)
	if left != right {
		sa.error(b.t, "Mismatched types for '=='. %v == %v", left, right)
	}
	return "bool" // doesnt matter left or right
}

func (b NotEqual) analyzeE(sa *SemAn) string {
	left := b.left.analyzeE(sa)
	right := b.right.analyzeE(sa)
	if left != right {
		sa.error(b.t, "Mismatched types for '!='. %v != %v", left, right)
	}
	return "bool" // doesnt matter left or right
}

func (b Less) analyzeE(sa *SemAn) string {
	left := b.left.analyzeE(sa)
	right := b.right.analyzeE(sa)
	if left != right {
		sa.error(b.t, "Mismatched types for '<'. %v < %v", left, right)
	} else if left != "int" && left != "float" { // doesnt matter left or right
		sa.error(b.t, "'<' needs int or float values, have %v and %v", left, right)
	}
	return "bool" // doesnt matter left or right
}

func (b Greater) analyzeE(sa *SemAn) string {
	left := b.left.analyzeE(sa)
	right := b.right.analyzeE(sa)
	if left != right {
		sa.error(b.t, "Mismatched types for '>'. %v > %v", left, right)
	} else if left != "int" && left != "float" { // doesnt matter left or right
		sa.error(b.t, "'>' needs int or float values, have %v and %v", left, right)
	}
	return "bool" // doesnt matter left or right
}

func (b LessEq) analyzeE(sa *SemAn) string {
	left := b.left.analyzeE(sa)
	right := b.right.analyzeE(sa)
	if left != right {
		sa.error(b.t, "Mismatched types for '<='. %v <= %v", left, right)
	} else if left != "int" && left != "float" { // doesnt matter left or right
		sa.error(b.t, "'<=' needs int or float values, have %v and %v", left, right)
	}
	return "bool" // doesnt matter left or right
}

func (b GreaterEq) analyzeE(sa *SemAn) string {
	left := b.left.analyzeE(sa)
	right := b.right.analyzeE(sa)
	if left != right {
		sa.error(b.t, "Mismatched types for '>='. %v >= %v", left, right)
	} else if left != "int" && left != "float" { // doesnt matter left or right
		sa.error(b.t, "'>=' needs int or float values, have %v and %v", left, right)
	}
	return "bool" // doesnt matter left or right
}

func (p Power) analyzeE(sa *SemAn) string {
	base := p.base.analyzeE(sa)
	power := p.base.analyzeE(sa)
	if power != "int" {
		sa.error(p.t, "Exponent for '^' must be int, have %v", power)
	} else if base != "int" && base != "float" {
		sa.error(p.t, "Base for '^' must be int or float, have %v", base)
	}
	return base
}

func (a And) analyzeE(sa *SemAn) string {
	right := a.right.analyzeE(sa)
	left := a.right.analyzeE(sa)
	if left != "bool" || right != "bool" {
		sa.error(a.t, "'&&' takes boolean values, have %v && %v", left, right)
	}
	return right // doesnt matter which one isnt a bool (if there is one), cause the
	// SemAn will exit
}

func (o Or) analyzeE(sa *SemAn) string {
	right := o.right.analyzeE(sa)
	left := o.right.analyzeE(sa)
	if left != "bool" || right != "bool" {
		sa.error(o.t, "'||' takes boolean values, have %v || %v", left, right)
	}
	return right // doesnt matter which one isnt a bool (if there is one), cause the
	// SemAn will exit
}

/**************************reciever goes back to sa***************************/

func (sa *SemAn) error(t token, msg string, args ...interface{}) {
	fullMsg := fmt.Sprintf(msg, args...)
	sa.errors = append(sa.errors, fmt.Sprintf("Line %v:", t.line)+fullMsg)
}

func (sa *SemAn) checkErrors() {
	if len(sa.errors) != 0 {
		fmt.Println("Semantic errors:")
		for _, s := range sa.errors {
			fmt.Println(s)
		}
		os.Exit(0)
	}
}
