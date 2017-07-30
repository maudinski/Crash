package go

import (
    "fmt"
    "strconv"
    "os"
)

type DeclaredFunc struct {
    name string // redundent cause it's being stored in a mao by it's name but ohwell
    params []string // slice of the types
    returnType string // will eventually be able to return multiple types
}

type SemAn struct {
    ast *Ast
    es EnvironmentStack
    errors []string
    declaredFuncs map[string]DeclaredFunc
    currentFunction string // holds the name of the current function being parsed
}                           // for analyzing return statments

func newSemAn(ast *Ast) *SemAn {
    sa := new(SemAn)
    sa.ast = ast
    sa.es = newEnvironemntStack()
    sa.errors = make([]string, 0)
    sa.funcs = make(map[string]DeclaredFunc, 0)
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
        sa.declaredFuncs[n] = DeclaredFunc{n, params, f.returnType})
    }

}

// calls analyze block on each function in the ast
func (sa *SemAn) phase2() {
    for name, f := range(sa.ast.functions) {
        sa.es.pushNewEnv() // new env created before each block is entered
        for _, p := f.params { // add all the parameters to the environement
            sa.es.add(p.id.value, p.ttype)
        }
        sa.currentFunction = name
        sa.analyzeBlock(f.block)
        sa.es.popEnv() // pop the environment
    }
}

func (sa *SemAn) analyzeBlock(block Block) {
    for _, statement := range(block){
        sa.analyzeStatementSwitch(statements)
    }
}

/******this switches control over to the statement********/

func (sa *SemAn) analyzeStatementSwitch(s Statement) {
    switch state := s.(type){
    default:
        state.analyze(sa)
    }
}

/********recievers from here on out turn to the statement and get sa passed**********/
/***turns out to a little less code in parse block, and eaiser modification**********/

// something like that
func (d Declaration) analyze(sa *SemAn) {
    exists, _ := sa.checkTop(d.id.value)
    if exists {
        sa.error(d.t, "Var %v already exists in this scope", d.id.value)
        return
    }
    sa.es.add(d.id.value, d.ttype)
    sa.analyzeExpression(d.ttype, d.value)
}

func (r Return) analyze(sa *SemAn) {
    if p.analyzeExpression(sa.declaredFuncs[sa.currentFunc].returnType, r.value)
}

func (i If) analyze(sa *SemAn) {
    sa.analyzeExpression("bool", i.exp)
    sa.es.pushNewEnv()
    sa.analyzeBlock(i.trueBlock)
    sa.es.popEnv()
    if !i.isElse{
        sa.es.pushNewEnv()
        sa.analyzeBlock(i.falseBlock)
        sa.popEnv()
    }
}

func (w While) analyze(sa *SemAn) {
    sa.analyzeExpression("bool", w.condition)
    sa.es.pushNewEnv()
    sa.analyzeBlock(w.block)
    sa.es.popEnv()
}

func (r Reassignment) analyze(sa *SemAn) {
    exists, ttype := sa.check(r.id.value)
    if !exists {
        sa.error(r.t, "Var %v not declared", r.id.value)
        return
    }
    sa.analyzeExpression(r.ttype, r.value)
}

// returns the return type, only needed if the function call is part of an expression
// still return either way
// seems right
func (c Call) analyze(sa *SemAn) string {
    declaredFunc := sa.declaredFuncs[sa.currentFunc]
    if len(declaredFunc.params) != len(c.params) {
        p.error(c.token, "Wrong number of arguments to function call %v", c.id.value)
    } else {
        for i, exp := range(c.params) {
            p.analzyeExpression(declaredFunc.params[i], exp)
        }
    }
    return declaredFunc.returnType
}

/*********this goes back to sa reciever but switches it over to expression*********/
// should validate the expression and make sure it is of the same type passed
// something like that // TODO TODO TODO remove that return at the beginning of the func
func (sa *SemAn) analyzeExpression(supposedType string, e Expression) {

                    return // TODO BUG BUG BUG BUG BUG BUG BUG BUG BUG BUG BUG BUG

    switch exp := e.(type){
    default: // pseudo code, just to show the concept. Need to check types
        ttype := exp.analyze() // these will return type, and
    }    // check it down here
    if ttype == "" {return} // it wasn't a valid expression. .analyze() will do the error
    if ttype != supposedType {
        p.error(e.t, "Expression %v not of correct type. Is %v, should be %v"
                                        , e.String(), ttype, supposedType)
    }
}

/*************reciever goes over to expression and takes sa as param*************/
// thinking these return the type of the expression
func (id Id) analyze(sa *SemAn) string {

}

func (c Call) analyze(sa *SemAn) string {

}

// and so on

/**************************reciever goes back to sa***************************/

func (sa *SemAn) error(t token, msg string, args ...interface{}) {
    fullMsg = fmt.Sprintf(msg, args...)
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