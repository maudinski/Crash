package main

import (
    "fmt"
    "os"
    "io/ioutil"
)

func main() {
    data, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        throwError("File not found")
    }
    d := newData(data)
    lexer := newLexer(d)
    lexer.setKeywords(/*"if", "while", */"print")
    lexer.setTypes("int", /*"float",*/ "string"/*, "bool"*/)

    for tok := range lexer.stream() {
        fmt.Println(tok)
        //fmt.Print(tok.value)
    }
}

func throwError(s ...interface{}) {
    fmt.Print("EXCEPTION THROWN: ")
    fmt.Println(s...)
    os.Exit(1)
}
