package main
// point of this is to check if variables exist in the scope, functions called where
// already defined/built in, expressions are valid, and type checking for expression
// and variable assignments/reassignments. Final string of error checking
import (
    "fmt"
    //"strconv" // atom doesn't recognize that comment lol
    "os"
)

type DeclaredFunc struct {
    name string // redundent cause it's being stored in a mao by it's name but ohwell
    params []string // slice of the types
    returnType string // will eventually be able to return multiple types
}

type SemAn struct {
    ast *Ast
    es *EnvStack
    errors []string
    // phase1. Holds info from function headers. Only for a more convenient format
    // than the []*Function in ast
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
func (sa *SemAn) analyze() {
    sa.phase1()
    sa.checkErrors() // will only be global variable expression errors
    sa.phase2()
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
    for _, dec := range(sa.ast.globals) { // add global env to stack, never popped
        dec.analyze(sa) // remeber that control switches to the statement, no confuse
    }   // analyze already adds it to the current environment
    for n, f := range(sa.ast.functions) { // name, *Function
        params := make([]string, 0) // get a []string of the parameter types
        for _, p := range(f.params) {
            params = append(params, p.ttype) // slice of the types
        }
        sa.declaredFuncs[n] = DeclaredFunc{n, params, f.returnType}
    }

}

// calls analyze block on each function in the ast
func (sa *SemAn) phase2() {
    for name, f := range(sa.ast.functions) {
        sa.es.pushNewEnv() // new env created before each block is entered
        for _, p := range(f.params) { // add all the parameters to the environement
            sa.es.add(p.id.value, p.ttype)
        }
        sa.currentFunc = name
        sa.analyzeBlock(f.block)
        sa.es.popEnv() // pop the environment
    }
}

/******this switches control over to the statement********/
func (sa *SemAn) analyzeBlock(block Block) {
    for _, s := range(block){
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
    if !i.isElse{
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

// TODO should not allow strings to be reassigned
func (r Reassignment) analyze(sa *SemAn) {
    exists, ttype := sa.es.check(r.id.value)
    if !exists {
        sa.error(r.t, "Var %v not declared", r.id.value)
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
    for i, exp := range(c.params) {
        sa.analyzeExpression(declaredFunc.params[i], exp)
    }
}

/*********this goes back to sa reciever but switches it over to expression*********/
// should validate the expression and make sure it is of the same type passed
// something like that // TODO TODO TODO remove that return at the beginning of the func
func (sa *SemAn) analyzeExpression(supposedType string, e Expression) {
                    return;         return
                        return;                return
            return;              return
                            return;         return
    // TODO putting this here so i remeber. Go make reassignment.analyze() not allow
    // strings to be reassigned, as it will complicate assembly code. So i guess that
    // means theyre 'immutable'
    // do it like this to do type checking
}

/*************reciever goes over to expression and takes sa as param*************/
// thinking these return the type of the expression
func (id Id) analyzeE(sa *SemAn) string {
    return ""
}
// will basically be c.analyze(sa) but with a return type
func (c Call) analyzeE(sa *SemAn) string {
    return ""
}
// and so on

/**************************reciever goes back to sa***************************/

func (sa *SemAn) error(t token, msg string, args ...interface{}) {
    fullMsg := fmt.Sprintf(msg, args...)
    sa.errors = append(sa.errors, fmt.Sprintf("Line %v:", t.line) + fullMsg)
}

func (sa *SemAn) checkErrors() {
    if len(sa.errors) != 0 {
        fmt.Println("Semantic errors:")
        for _, s := range(sa.errors) {
            fmt.Println(s)
        }
        os.Exit(0)
    }
}