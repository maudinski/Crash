package go

import (
    "fmt"
    "strconv"
    "os"
)

type DeclaredFunc struct {
    name string
    params []string
    returnType string // will eventually be able to return multiple types
}

type SemAn struct {
    ast *Ast
    es EnvironmentStack
    errors []string
    funcs []DeclaredFunc
}

func newSemAn(ast *Ast) *SemAn {
    sa := new(SemAn)
    sa.ast = ast
    sa.es = newEnvironemntStack()
    sa.errors = make([]string, 0)
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
// should be good
func (sa *SemAn) phase1() {
    sa.es.pushNewEnv()
    for _, dec := range(sa.ast.globals) { // add global env to stack, never popped
        dec.analyze(sa) // remeber that control switches to the statement, no confuse
    }   // analyze already adds it to the current environment
    for n, f := range(sa.ast.functions) { // name, *Function
        params := make([]string, 0) // get a []string of the parameter types
        for _, p := range(f.params) {
            params = append(params, p.ttype)
        }
        sa.funcs = append(sa.funcs, DeclaredFunc{n, params, f.returnType})
    }

}
// just loop through functions
//
// Gather up called function names and number of parameters passed, check on
// already gathered functions
//
// set up environments for each function block, pseduo environments for if's and
// whiles, etc. Environements store variables and their types
//
// check for valid expressions based on globals and environment variables
// valid expression just means that any variables and functions used is an expression have
// already been declared with correct type/ return correct type (respectively)
func (sa *SemAn) phase2() {
    for _, f := range(sa.ast.functions) {
        sa.es.pushNewEnv() // new env created before each block is entered
        for _, p := f.params { // add all the parameters to the environement
            sa.es.add(p.id.value, p.ttype)
        }
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

func (c Call) analyze(sa *SemAn) {
    
}

/*********this goes back to sa reciever but switches it over to expression*********/
// logic here is you need to check it against the type it's supposed to be
func (sa *SemAn) analyzeExpression(supposedType string, e Expression) {
    return // TODO
    switch exp := e.(type){
    default: // pseudo code, just to show the concept. Need to check types
        exp.analyze()
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

func (sa *SemAn) error(t token, msg string) {
    sa.errors = append(sa.errors, fmt.Sprintf("Line %v:", t.line) + msg)
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