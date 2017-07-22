// if lexer uses it, but it doesn't have lexer as a reciever, it's in lexerSecondaries.go
package main

import (
    "regexp"
    "fmt"
    "os"
)

 type Lexer struct {
     data *Data
     keyWords []string
     types []string
     lineNum int // tracked by new lines
     queued *Queue
     errors []string
}

func newLexer(data *Data) *Lexer {
    lx := new(Lexer)
    lx.data = data
    lx.lineNum = 1
    lx.queued = newQueue() // in lexerSecondaries.go
    lx.errors = make([]string, 0)
    return lx
}

func (lx *Lexer) setKeywords(kw ...string) {
    lx.keyWords = kw
}

func (lx *Lexer) setTypes(t ...string) {
    lx.types = t
}
// recursive for tabs and white space-
//this is the head. control is usually switched over to a second hand function
//once the beginning character or a token is identified. if the char is " " or "#"
//(blank or comment), it will discard approriately(spelling) and recursively call its
//self
//also, all the error checking (as of now) is done in the 'second hand' functions.
//the error strings are added to a the lx.errors. On the last call to this frunction
//from the parser (it will be the last call because this function will be returning false)
//, if the errors []string isnt empty, it will print them all out and exit the program
func (lx *Lexer) next() token {
    if !lx.queued.isEmpty() {
        return lx.queued.pull()
    }
    c := lx.data.next()
    if c == "@" { // arbitrary, i chose this to mean EOF. Will be fine within strings
        if len(lx.errors) != 0 { // reached only if nextByte is out of data
            for _, s := range(lx.errors){ fmt.Println(s) }
            os.Exit(0)
        }
        return token{"@", "@"}
    }
    if c == " " || c == "\t" { return lx.next() // spaces
    } else if c == "\n" { lx.lineNum++; return token{"NEWLINE", "\\n"} // \n
    } else if isWrapper(c) { return token{c, c} // any parenthesis
    } else if c == "#" { lx.scrapeComment(); return lx.next()  // comments
    } else if c == "." { return token{"DOT_OP", c} // dot operator, idk what to do
    } else if isDigit(c) { return lx.getNumToken(c) // ints and floats
    } else if isOperator(c) { return lx.getOpToken(c) // math and bool operators
    } else if c == "\"" { return lx.getStrToken() // strings
    } else { return lx.getAmbiguousToken(c) } // cant classify by single byte
}

func (lx *Lexer) goBack(t token) {
    lx.queued.push(t)
}
//maybe
func (lx *Lexer) peek() token {
    t := lx.next()
    lx.goBack(t)
    return t
}

func (lx *Lexer) scrapeComment() {
    for c := lx.data.next(); c != "@" && c != "\n"; c = lx.data.next() {}
    lx.lineNum++
}

//new lines arent allow in strings
func (lx *Lexer) getStrToken() token {
    str := ""
    closed := false
    for c := lx.data.next(); c != "@"; c = lx.data.next(){
        if c == "\n" {
            closed = false
            lx.lineNum++
            break
        } else if c == "\"" {
            closed = true
            break
        }
        str += c
    }
    if !closed {
        lx.errors=append(lx.errors, "String never closed, line " + toString(lx.lineNum-1))
    }
    return token{"STRING", str}
}

// might have to have specified operator types, like BOOL_OPERATOR and MATH_OPERATOR
func (lx *Lexer) getOpToken(op string) token {
    c := lx.data.peek()
    if c == "@" { // more of a parsing error than lexing, but ehh, we're here
        lx.errors = append(lx.errors, "Nothing after operator, line "+toString(lx.lineNum))
    } else if isWrapper(c) {
        lx.data.next()
        op += c
    }
    return token{"OPERATOR", op}
}

// checks to see if the number ended at an approriate spot(to the extent o the lexers
// knowledge)
// does not differentiate between int and float, but does account for both. returns
//"NUMBER" as type of token regardless
func (lx *Lexer) getNumToken(snum string) token {
    c := ""
    for c = lx.data.next(); c != "@"; c = lx.data.next() {
        if !isDigit(c) { break }
        snum += c
    }
    if !canNumEndHere(c) {
        lx.errors = append(lx.errors, "Invalid number on line " + toString(lx.lineNum))
    }
    lx.data.goBack() // since one extra byte was grabbed

    return token{"NUMBER", snum}
}

// str is the first character of the token (already read in by lx.next())
func (lx *Lexer) getAmbiguousToken(str string) token {
    for c := lx.data.next(); c != "@"; c = lx.data.next() {
        if ok, _ := regexp.MatchString("[_a-zA-Z0-9]", c); !ok { break }
        str += c
    }
    lx.data.goBack() // since one extra was read in
    // now classify the ambiguous token
    if lx.isType(str) {
        return token{"TYPE", str}
    } else if lx.isKeyword(str) {
        return token{"KEYWORD", str}
    }
    return token{"ID", str}
}

// only used by getAmbiguousToken()
func (lx *Lexer) isType(s string ) bool {
    for _, ttype := range lx.types { if ttype == s { return true } }
    return false
}

// only used by getAmbiguousToken()
func (lx *Lexer) isKeyword(s string ) bool {
    for _, kw := range lx.keyWords { if kw == s { return true } }
    return false
}


















