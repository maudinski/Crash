package main

import (
    "fmt"
)

type node interface{
    exec()
}

type Declaration struct {
    ttype string
    name string
    value string
}

func (d Declaration) exec() {

}

type Print struct {
    value string
}

func (p Print) exec() {
    fmt.Println(p.value)
}