package main

import (

)

type data struct {
    data []byte
    index int
    length int
}

func newData (fileData []byte) *data {
    d := new(data)
    d.data = fileData
    d.length = len(fileData)
    d.index = -1 //-1 cause it starts with an incriment
    return d
}

func (d *data) nextByte() (byte, bool){
    d.index++
    if d.index == d.length {
        return ' ', false
    }
    return d.data[d.index], true
}

func (d *data) peek() (byte, bool) {
    if d.index + 1 == d.length {
        return ' ', false
    }
    return d.data[d.index+1], true
}

func (d *data) goBack(){
    d.index--
}
