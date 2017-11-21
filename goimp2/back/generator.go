make the symbol table. todo in semAn for notes
// will output the generated code to a file with fileName.asm
package main

import (
      "os"
)

type Generator struct {
      ast *Ast
      // file name
      fn string
      // file for output
      file *os.File
      // symbol table
      symbols map[string]int // not sure
}

func newCodeGenerator(ast *Ast, fn string) *Generator{
      gen := new(Generator)
      gen.ast = ast
      gen.fn = fn
      gen.symbols = make() //something
      return gen
}

// brain
func (gen *Generator) generate() {
      // open the file
      gen.file, err := os.Create(g.fileName+".asm")
      if err != nil {
            fmt.Println("Error creating assembly file:", err)
            os.Exit(0)
      }
      defer gen.file.Close()
      // starts generating, storing activation record templates, idk yet
      // when compiling main function, toss in a mov rax, 60; mov, rdi 0; syscall
      // to exit. Put main function at the top of the code generated, and everything
      // else beneth it. All other functions will have returns at the end of their
      // prcodures
}

// writes to the file, exits if error
func (gen *Generator) write(ass *string) {
      _, err := gen.file.Write(*ass)
      if err != nil {
            fmt.Println("Error writing to assembly file:", err)
            os.Exit(0)
      }
}



























