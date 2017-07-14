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

type Declaration struct {
    n Name
    t Type
    v Value
}

type Print struct {
    params []node
}
func newPrint(ns ...node) Print {
    p := Print{ns}
    return p
}
























