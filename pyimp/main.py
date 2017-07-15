import sys
import stream
import lexer
import parser
import AST

def main():
    f = open(sys.argv[1] ,"r")
    filestream = stream.Stream(f.read())
    keywords = ["print", "if"]
    types = ["int", "float", "string", "bool"]
    lx = lexer.Lexer(filestream, keywords, types)
    p = parser.Parser(lx)
    ast = AST.Ast(p)
    ast.interpret()


main()