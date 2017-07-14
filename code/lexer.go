package main

import(
    "errors"
    "regexp"
)

type lexer struct {
    data *data
    keyWords []string
    types []string
    newLines int //count of new lines
    stream chan token
}

func newLexer(d *data) *lexer {
    l := new(lexer)
    l.data = d
    l.newLines = 0
    return l
}

func (lx *lexer) setKeywords(kw ...string) {
    lx.keyWords = kw
}

func (lx *lexer) setTypes(t ...string) {
    lx.types = t
}

func (lx *lexer) isType(s string ) bool {
    for _, ttype := range lx.types { if ttype == s { return true } }
    return false
}

func (lx *lexer) isKeyword(s string ) bool {
    for _, kw := range lx.keyWords { if kw == s { return true } }
    return false
}


func (lx *lexer) flow() chan token {
    lx.stream = make(chan token)
    go func() {
        for b, ok := lx.data.nextByte(); ok; b, ok = lx.data.nextByte() {
            c := string(b)
            if c == " " {
                continue
            } else if c =="\n"{
                lx.newLines++//maybe?
                lx.stream <- token{"NEWLINE", "\\n"}//BUG may cause bug cause escape seq
            } else if isWrapper(c) {//if its some sort of parenthesis/bracket
                lx.stream <- token{c, c}
            } else if isSingleOperator(c) {//if its an operator
                b2, ok := lx.data.peek() //TODO this could be its own function
                c2 := string(b2)
                if ok && isSingleOperator(c2) {
                    lx.data.nextByte()
                    c += c2
                }
                lx.stream <- token{c, c}
            } else if c == "\"" {//if its about to be a string
                str, err := lx.getTill("\"")
                if err != nil { throwError("String never closed") } // only error check
                lx.stream <- token{"STRING", str}
            } else if c == "#" { //if its a comment
                lx.getTill("\n")
            } else if isDigit(c) {
                lx.stream <- token{"NUMBER", lx.getNumber(c)}
            } else {
                word := lx.getWord(c)
                if lx.isType(word) {
                    lx.stream <- token{"TYPE", word}
                } else if word == "print" {
                    lx.stream <- token{"PRINT", word}
                } else if lx.isKeyword(word){
                    lx.stream <- token{"KEYWORD", word}
                } else {
                    lx.stream <- token{"VARIABLE", word}
                }
            }
        }
        close(lx.stream)
    }()
    return lx.stream
}

//TODO put these in lx so not constant compile (the r's)

func (lx *lexer) getWord(firstLetter string) string {
    r, _ := regexp.Compile("[_a-zA-Z0-9]")
    word := firstLetter
    for b, ok := lx.data.nextByte(); ok; b, ok = lx.data.nextByte() {
        c := string(b)
        if !r.MatchString(c){ break }
        word += c
    }
    lx.data.goBack()
    return word
}

func (lx *lexer) getTill(stop string) (string, error) {
    str := ""
    for b, ok := lx.data.nextByte(); ok; b, ok = lx.data.nextByte(){
        c := string(b)
        if c == stop {
            return str, nil
        }
        str += c
    }
    return "", errors.New("")
}

func (lx *lexer) getNumber(firstDigit string) string {
    snum := firstDigit
    var c string
    for b, ok := lx.data.nextByte(); ok; b, ok = lx.data.nextByte() {
        c = string(b)
        if !isDigit(c) { break }
        snum += c
    }
    //IDEA call canNumberEndHere to early check
    lx.data.goBack()
    return snum
}

func isDigit(c string) bool {
    r, _ := regexp.Compile("[.0-9]")
    return r.MatchString(c)
}

func canNumberEndHere(c string) bool {
    //BUG over all this function isnt complete
    return isWrapperRight(c) || isSingleOperator(c) || c=="\n" || c=="#" || c==" "
}

func isWrapper(s string) bool {
    r, _ := regexp.Compile("[{}()\\[\\]]")
    return r.MatchString(s)
}

func isSingleOperator(s string) bool {
    r, _ := regexp.Compile("[=\\-\\+/%><!]")
    return r.MatchString(s)
}

func isWrapperRight(s string) bool {
    return s == "}" || s == "]" || s == ")"
}
