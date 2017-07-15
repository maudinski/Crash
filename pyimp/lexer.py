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
                return ("STRING", self.get_matched("[^\"]", "", False)) # what if no " ?
            elif c == "#":
                self.get_matched("[^\n]", c, True)
                pass
            elif re.match("[.0-9]", c):
                return ("NUMBER", self.get_matched("[.0-9]", c, True))
            elif c == ",":
                return (c, c)
            else:
                word = self.get_matched("[_a-zA-Z0-9]", c, True)
                if word in self.keywords:
                    return (word, word)
                elif word in self.types:
                    return ("TYPE", word)
                else:
                    return ("NAME", word)

            c = self.stream.next()
        return None

    def put_back(self, token):
        self.backedUp.insert(0, token)

    def peek(self):
        tok = self.next()
        self.put_back(tok)
        return tok

    def get_matched(self, regex, s, put_first_nonmatch_back):
        c = self.stream.next()
        while re.match(regex, c):
            s += c
            c = self.stream.next()
        if put_first_nonmatch_back:
            self.stream.go_back()
        return s










