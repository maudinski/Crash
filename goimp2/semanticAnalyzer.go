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

// pre shit
// 1) pulls globals from sa.ast.globals
// 2) pulls defined functions from ast.functions, they're parameter types in order,
//    and the return type(s) in more accessible manner
// phase really makes this sound way more epic than it is
// should be good
func (sa *SemAn) phase1() {
    sa.es.pushNewEnv()
    for _, dec := range(sa.ast.globals) {
        sa.es.add(dec.id.value, dec.ttype)
    }
    for n, f := range(sa.ast.functions) { // name, *Function
        params := make([]string, 0) // get a []string of the parameter types
        for _, p := range(f.params) {
            params = append(params, p.ttype)
        }
        sa.funcs = append(sa.funcs, DeclaredFunc{n, params, f.returnType})
    }
}

// Gather up called function names and number of parameters passed, check on
// already gathered functions
//
// set up environments for each function block, pseduo environments for if's and
// whiles, etc. Environements store variables and their types
//
// check for valid expressions based on globals and environment variables
//
//
func (sa *SemAn) phase2() {
    // when entering a function, you gotta add the parameters to the environment
    // as already being declared
}












