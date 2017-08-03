package main

import (
	"fmt"
	"regexp"
	"strconv"
)

type token struct {
	ttype string
	value string
	line  int
}

func (t token) String() string {
	return fmt.Sprintf("Type: %-15v Value: %-15v Line: %-15v", t.ttype, t.value, t.line)
}

//Queue implimentation for the returned or peeked tokens. In relation to whatever the
//parser does with the lexer
type Queue struct {
	tokens []token
	amt    int
}

func newQueue() *Queue {
	q := new(Queue)
	q.tokens = make([]token, 0)
	q.amt = 0
	return q
}

func (q *Queue) push(t token) {
	q.tokens = append(q.tokens, t)
	q.amt++
}

func (q *Queue) pull() token {
	t := q.tokens[0]
	if q.amt == 1 {
		q.tokens = make([]token, 0) // idk about this, might be allocating mem
	} else {
		q.tokens = q.tokens[1:]
	}
	q.amt--
	return t
}

/***/
func (q *Queue) isEmpty() bool {
	return q.amt == 0
}

func isWrapper(s string) bool {
	ok, _ := regexp.MatchString("[{}()\\[\\]]", s)
	return ok
}

func isOperator(s string) bool {
	ok, _ := regexp.MatchString("[=\\-\\+/%>*<!]", s)
	return ok
}

func isDigit(c string) bool {
	ok, _ := regexp.MatchString("[.0-9]", c)
	return ok
}

// used by getNumToken to check if the byte right after the last digit in the num is
// something valid that a number could be attached to (so an operator, a blank space, etc)
func canNumEndHere(c string) bool {
	if ok, _ := regexp.MatchString("[ ,;\\])\n]", c); !ok {
		if c != "EOF" {
			return isOperator(c)
		}
	}
	return true
}

// too long for me
func toString(i int) string {
	return strconv.Itoa(i)
}
