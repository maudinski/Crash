package main

import (

)

type EnvStack struct {
    stack []Environment
    farthestPos int
}

func newEnvironmentStack() *EnvStack {
    es := new(EnvStack)
    es.stack = make([]Environment, 0)
    es.farthestPos = -1 // in accordance to array indexing
    return es
}
// maybe
func (es *EnvStack) pushNewEnv() {
    if len(es.stack) > es.farthestPos+1 {
        es.stack[es.farthestPos + 1] = newEnvironment()
    } else {
        es.stack = append(es.stack, newEnvironment())
    }
    es.farthestPos++
}

func (es *EnvStack) popEnv() {
    es.stack[es.farthestPos] = nil
    es.farthestPos--
}

// adds to the top most environment
func (es *EnvStack) add(id string, ttype string) {
    es.stack[es.farthestPos].add(id, ttype)
}

// checks all environments for existing var name
// seems alright
func (es *EnvStack) check(id string) (bool, ttype) {
    for i := 0; i <= es.farthestPos; i++ {
        ok, ttype := es.stack[i].check(id)
        if ok {
            return true, ttype
        }
    }
    return false, ""
}

/******/
type Environment struct {
    vars map[string]string // map[id] gives variables type
}

func newEnvironment() *Environment {
    env := new(Environement)
    env.vars = make(map[string]string, 0) // i think
}

func (env *Environment) add(id string, ttype string) {
    env.vars[id] = ttype
}

func (env *Environment) check(id string) (bool, string) {
    ttype, ok := env.vars[id]
    return ok, ttype
}