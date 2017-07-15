package main

import(

)

type token struct {
    ttype string
    value string
}

const NEWLINE = 0
const LBRACKET = 1
const RBRACKET = 2
const LPAREN = 3
const RPAREN = 4

const PRINT = 7
const STRING = 8
