package main

import (
    "fmt"
    "os"
    "io/ioutil"
)


func main() {
    if len(os.Args) != 2 {
        fmt.Println("Enter one file")
        os.Exit(1)
    }
    fileData, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        fmt.Println("File not read: ", os.Args[1])
        os.Exit(1)
    }
    data := newData(fileData, os.Args[1])
    lexer := newLexer(data)
    lexer.setKeywords("if", "func", "while", "for", "return",
                        "struct", "global")
    lexer.setTypes("int", "float", "string", "char", "byte")
    for t := lexer.next(); t.ttype != "EOF"; t = lexer.next(){
        fmt.Println(t)
    }
    //parser = newParser(lexer)
    //ast = parser.parse()
}
