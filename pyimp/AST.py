import e
import n
import sys

class Ast:

    def __init__(self, parser):
        self.p = parser
        self.vars = {}

    def interpret(self):
        node = self.p.next()
        while node is not None:
            self.exec_switch(node)
            node = self.p.next()

    # only return if the calling function could expect a return type
    def exec_switch(self, node):
        if type(node) == n.Print:
            return self.exec_print(node)
        if type(node) == n.Declaration:
            return self.exec_declaration(node)
        if type(node) == n.Number or type(node) == n.String:
            return self.exec_basic(node)
        if type(node) == n.Name:
            return self.exec_name(node)
        if type(node) == n.Operation:
            return self.exec_operation(node)

    def exec_print(self, p_n):
        for node in p_n.params:
            if type(node) == n.Name:
                print(self.exec_switch(node)[1], "",end='')# name gives type and val
            else:
                print(self.exec_switch(node), "", end='')
        print()
        return None

    # for nodes that just hold a string value (Number and String)
    def exec_basic(self, node):
        return node.value

    def exec_declaration(self, d_n):
        self.vars[d_n.name.value] = [d_n.ttype, self.exec_switch(d_n.value)]

    # if this is called it because the VALUE for name is needed, not because
    # it needs to be added. Thats exec_declaration
    def exec_name(self, n_n):
        if n_n.value not in self.vars:
            e.throwError("{} undefined on line ".format(n_n.value), self.p.line)
        return self.vars[n_n.value]

    def exec_operation(self, o_n):
        op = o_n.operator
        left = self.exec_switch(o_n.left)
        right = self.exec_switch(o_n.right)
        if e.strictlyString(left) or e.strictlyString(right):
            if op != "+":
                e.throwError("String operations only support '+' on line ", self.p.line)
        elif e.isFloat(left) or e.isFloat(right):
            left, right = float(left), float(right)
        else:
            left, right = int(left), int(right)
        if op == "+": return left + right
        if op == "-": return left - right
        if op == "/": return left / right
        if op == "%": return left % right
        if op == "*": return left * right








