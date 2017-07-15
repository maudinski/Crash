# stands for nodes, too long to type it out all the time

class Print:
    def __init__(self, params):
        self.params = params
    def to_string(self):
        s = ""
        for p in self.params:
            s += p.to_string() + " "
        return "Print: "+s

class Operation:
    def __init__(self, operator, left, right):
        self.operator = operator # string of the operator
        self.left = left # Number or String node
        self.right = right # Number or String ode
    def to_string(self):
        return " Operation: "+self.left.to_string() + self.operator + self.right.to_string()

class Number:
    def  __init__(self, value):
        self.value = value # a string still
    def to_string(self):
        return " Number: "+self.value

class String:
    def __init__(self, value):
        self.value = value # a python string
    def to_string(self):
        return " String: "+self.value

class Name:
    def __init__(self, value):
        self.value = value# prolly a string
    def to_string(self):
        return " Name: "+self.value

class Declaration:
    def __init__(self, ttype, name, value):
        self.ttype = ttype # a string
        self.name = name # a Name object
        self.value = value # a String, Number, or Operation(eventuall bool or func too)
    def to_string(self):
        return "Declaration: "+self.ttype +' '+self.name.to_string()+' '+self.value.to_string()















