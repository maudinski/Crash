package main

import ("strconv")

func (p *Parser) parseExpression() Expression {
    stack := newTokenStack()
    postfix := make([]ExpressionPart, 0)
    isBool := false // if a boolean operator is reached set this true
    isInt := false // if any operator besides boolean and + is reached set this to true
    var t token
    for t = p.lx.next(); isExpressionPart(t); t = p.lx.next() {
        if isOpOrParen(){
            // evaluate the stack, potentially add to postfix
        } else {
            postfix = append(postfix, /*some approriate node*/)
        }
    }
    for !stack.isEmpty() {
        t := stack.pop()
        if t.value == "(" {
            // theres an error
        } else {
            postfix = append(postfix, /*some approriate node*/)
        }
    }
    return Postfix{t: t, exp: postfix, overallType: somebullshit}
}

type TokenStack struct {
    stack []token
    amt int
}

func newTokenStack() *TokenStack {
    ts := new(TokenStack)
    ts.stack = make([]token, 0)
    ts.amt = 0
    return ts
}

func (ts *TokenStack) isEmpty() bool {
    return ts.amt == 0
}

// look over
func (ts *TokenStack) push(t token) {
    if len(ts.stack) == ts.amt {
        ts.stack = append(ts.stack, t)
    } else {
        ts.stack[ts.amt] = t
    }
    ts.amt++
}

func (ts *TokenStack) pop() token {
    if es.isEmpty {
        return token{"EMPTY", "EMPTY"}
    }
    t := ts.stack[ts.amt-1]
    ts.amt--
    return t
}