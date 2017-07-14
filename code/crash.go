package main

import (
    "fmt"
    "os"
    "io/ioutil"
)

func main() {
    fd, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        throwError("File not found")
    }
    data := newData(fd)

    lexer := newLexer(data)
    lexer.setKeywords(/*"if", "while", */"print")
    lexer.setTypes("int", /*"float",*/ "string"/*, "bool"*/)

    parser := newParser(lexer)
    ast := parser.parse()
    ast.interpret()
}

func throwError(s ...interface{}) {
    fmt.Print("EXCEPTION THROWN: ")
    fmt.Println(s...)
    os.Exit(1)
}
