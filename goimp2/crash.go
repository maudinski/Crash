package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

var outputAssembly = false

func main() {
	var fileName string
	if len(os.Args) == 1 {
		fmt.Println("No file specified")
		os.Exit(0)
	}
	if os.Args[1] == "-ass" {
		outputAssembly = true
		fileName = os.Args[2]
	} else {
		fileName := os.Args[1]
	}
	fileName = checkFileName(fileName)
	fileData, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("File not read: '%v'\n", fileName)
		os.Exit(1)
	}
	data := newData(fileData, fileName)
	lexer := newLexer(data)
	lexer.setKeywords("if", "func", "while" /*"for",*/, "return",
		"struct", "global")
	lexer.setTypes("int", "float", "string", "bool")
	/*for t := lexer.next(); t.value != "EOF"; t = lexer.next() {
		fmt.Println(t)
	}*/
	parser := newParser(lexer)
	parser.setPrattMaps() // hard coded in ___ExpressionParser.go
	ast := parser.parse()
	// fmt.Println(ast)
	analyzer := newSemAn(ast)
	analyzer.analyze()
	fmt.Println(ast)
	//generator := newCodeGenerator(ast, symtbl, fileName)
	// outputs the generated code to a file with fileName.asm
	//generator.generate()
	//processAss(fileName)
}

func processAss(fileName string) {
	if outputAssembly {
		os.Exit(1)
	}
	// run ___embler on fileName.asm
	// run the linker
	// delete the ___embly file and .o file
}

// checks the file extension and truncates it
func checkFileName(fn string) string {

}












