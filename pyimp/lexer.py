import re
# gonna have to change shit up when i impliment arrays, specifically with the "word"
# else statment
class Lexer():
    def __init__(self, stream, keywords, types):
        self.stream = stream
        self.keywords = keywords
        self.types = types
        self.backedUp = []

    def next(self):
        if self.backedUp != []:
            return self.backedUp.pop()
        c = self.stream.next()
        while c is not None:
            if c == " ":
                pass
            elif c == "\n":
                return ("NEWLINE", "\\n")
            elif c in "{[()]}":
                return (c, c)
            elif c in "+-/%*":
                return ("OPERATOR", c)
            elif c in "<>=":
                if self.stream.peek() == "=":
                    c += self.stream.next()
                return (c, c)
            elif c == "\"":
                return ("STRING", self.get_string())# what if no " ?
            elif c == "#":
                self.dispose_comment()
                pass
            elif re.match("[.0-9]", c):
                return ("NUMBER", self.get_num(c))
            elif c == ",":
                return (c, c)
            else:
                word = self.get_word(c)
                if word in self.keywords:
                    return (word, word)
                elif word in self.types:
                    return ("TYPE", word)
                else:
                    return ("NAME", word)

            c = self.stream.next()
        return None

    def get_word(self, s):
        c = self.stream.next()
        while re.match("[_a-zA-Z0-9]", c): #BUG
            s += c
            c = self.stream.next()
        self.stream.go_back()
        return s

    def get_num(self, n):
        c = self.stream.next()
        while re.match("[.0-9]", c): #re would be best
            n += c
            c = self.stream.next()
        self.stream.go_back()
        return n


    def dispose_comment(self):
        c = self.stream.next()
        while c != "\n":
            c = self.stream.next()


    def get_string(self):
        returner = ""
        c = self.stream.next()
        while c != "\"":
            returner += c
            c = self.stream.next()
        return returner

    def put_back(self, token):
        self.backedUp.insert(0, token)

    def peek(self):
        tok = self.next()
        self.put_back(tok)
        return tok











