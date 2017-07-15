import sys
import stream
import lexer

def main():
    file = open(sys.argv[3] ,"r")
    st = stream.Stream(file.read())
    keywords = ["print", "if"]
    types = ["int", "float", "string", "bool"]
    for token in lexer.lex(st, keywords, types):
        print(token)

main()