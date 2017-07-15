import n
import error

class Parser:
    def __init__(self, lx):
        self.lx = lx
        self.line = 1

    def parse(self):
        token = self.lx.next()
        ast = []
        while token is not None:
            if token[0] == "NEWLINE":
                line += 1
            elif token[0] == "TYPE":
                ast.append(self.get_declarationN(token))
            elif token[0] == "print":
                ast.append(self.get_printN())
            else:
                print("(parse)not doing anything with", token)


    def get_declarationN(self, token):


    # janky but when those functions are wwritten, this should be able to do
    # print(2, "hello", x + 3)
    def get_printN():
        tok = self.lx.next()
        params = []
        if tok[0] != "(":
            error.throwError("expecting '(' after print on line ", self.line)
        while tok is not None:
            peekedTok = self.lx.peek()
            if peekedTok[0] == "OPERATOR":
                params.append(self.get_operation(tok))# gotta write
            elif tok[0] == "NAME":
                    params.append(n.Name(tok[1]))# gotta write
            elif tok[0] == "NUMBER":
                    params.append(n.Number(tok[1]))# gotta write
            elif tok[0] == "STRING":
                    params.append(n.String(tok[1]))# gotta write

            tok = self.lx.next()
            if tok[0] != ")" and tok[0] != ",":
                error.throwError("expecting ')' or ',' in print statement on line ", self.line)
            elif tok[0] == ")":
                break

        if params == []:
            error.throwError("no paramaters in print on line ", self.line)
        return n.Print(params)







