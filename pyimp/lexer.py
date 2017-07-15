def lex(stream, keywords, types):
    c = stream.next()
    while c is not None:
        if c == " ":
            pass
        elif c == "\n":
            yield ("NEWLINE", "\\n")
        elif c in "{[()]}":
            yield (c, c)
        elif c in "+-/%*":
            yield ("OPERATOR", c)
        elif c in "<>=":
            if stream.peek() == "=":
                c += stream.next()
            yield (c, c)
        elif c == "\"":
            yield ("STRING", get_string(stream))
        elif c == "#":
            dispose_comment(stream)
            pass
        elif c in ".0123456789": #re would be best
            yield ("NUMBER", get_num(c, stream))
        else:
            word = get_word(c, stream)
            if word in keywords:
                yield (word, word)
            elif word in types:
                yield ("TYPE", word)
            else:
                yield ("NAME", word)

        c = stream.next()


def get_word(s, stream):
    c = stream.next()
    while c in "abcdefghijklmnopqrstuvwxyz": #BUG
        s += c
        c = stream.next()
    stream.go_back()
    return s

def get_num(n, stream):
    c = stream.next()
    while c in ".0123456789": #re would be best
        n += c
    stream.go_back()
    return n


def dispose_comment(stream):
    c = stream.next()
    while c != "\n":
        c = stream.next()


def get_string(stream):
    returner = ""
    c = stream.next()
    while c != "\"":
        returner += c
        c = stream.next()
    return returner














