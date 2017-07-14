package main

import (
)

type node interface{}

type Name struct {
    value string
}

type Type struct {
    value string
}
//used for declaration
type Value struct {
    value string
}

type Reassign struct {
    n Name
    v Value
}

type Declaration struct {
    n Name
    t Type
    v node
}

type Print struct {
    params []node
}
func newPrint(ns ...node) Print {
    p := Print{ns}
    return p
}

type Operation struct {
    op Operator
    vLeft node
    vRight node
}

type Operator struct {
    v Value
}



















