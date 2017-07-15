import sys
import stream
import lexer
import parser

def main():
    f = open(sys.argv[1] ,"r")
    filestream = stream.Stream(f.read())
    keywords = ["print", "if"]
    types = ["int", "float", "string", "bool"]
    lx = lexer.Lexer(filestream, keywords, types)
    ast = parser.parseLx(lx)

main()