import n
import e

class Parser:
    def __init__(self, lx):
        self.lx = lx
        self.line = 1

    def next(self):
        token = self.lx.next()
        while token is not None:
            if token[0] == "NEWLINE":
                self.line += 1
            elif token[0] == "TYPE":
                return self.get_declarationN(token)
            elif token[0] == "print":
                return self.get_printN()
            else:
                print("(parse)not doing anything with", token)
            token = self.lx.next()
        return None

# pretty lengthy ugly functions cause they're actually error checking
    def get_declarationN(self, ttype):
        name = self.lx.next()
        equal = self.lx.next()
        value = self.lx.next()
        if name[0] != "NAME":
            e.throwError("Need variable name after '{}' on line ".format(ttype[1]), self.line)
        if equal[0] != "=":
            e.throwError("Need = after '{}' on line".format(name[1]), self.line)
        dec_val_n = ""
        if value[0] == "STRING":
            dec_val_n = n.String(value[1])
        elif value[0] == "NUMBER":
            dec_val_n = n.Number(value[1])
        node = dec_val_n#TODO comment this out once this works
        if self.lx.peek()[0] == "OPERATOR":
            dec_val_n = self.get_operation(node)

        return n.Declaration(ttype[1], n.Name(name[1]), dec_val_n)

    # janky but when those functions are wwritten, this should be able to do
    # print(2, "hello", x + 3)
    # BUG the logic here might be wrong
    # this may be pretty similar to what get_function or someshit will have to do, figure
    # out later tho
    def get_printN(self):
        tok = self.lx.next()
        params = []
        node = ""
        if tok[0] != "(":
            e.throwError("expecting '(' after print on line ", self.line)
        while tok is not None:
            # this chunk can prolly be factored out of this, Declaration, and Reassign
            tok = self.lx.next()
            if tok[0] == "NAME":
                node = n.Name(tok[1])
            elif tok[0] == "NUMBER":
                node = n.Number(tok[1])
            elif tok[0] == "STRING":
                node = n.String(tok[1])
            else:
                e.throwError("Invalid print parameter '{}' on line ".format(tok.value), self.line)
            peekedTok = self.lx.peek()
            if peekedTok[0] == "OPERATOR":
                params.append(self.get_operation(node))
            else:
                params.append(node)

            tok = self.lx.next()
            if tok[0] != ")" and tok[0] != ",":
                e.throwError("expecting ')' or ',' in print statement on line ", self.line)
            elif tok[0] == ")":
                break

        return n.Print(params)

    # takes an already created node for the left of the operator
    # operations can be of different types, they just always adhere to the more specific
    # type, so float + int is float, float + string is string. Bools would be strings in
    # all operations
    def get_operation(self, left_n):
        op_t = self.lx.next()
        right_t = self.lx.next()
        if right_t[0] == "NAME":
            return n.Operation(op_t[1], left_n, n.Name(right_t[1]))
        if right_t[0] == "STRING":
            return n.Operation(op_t[1], left_n, n.String(right_t[1]))
        if right_t[0] == "NUMBER":
            return n.Operation(op_t[1], left_n, n.Number(right_t[1]))
        e.throwError("Invalid operation on line ", self.line)





