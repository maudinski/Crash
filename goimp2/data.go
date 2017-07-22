package main

import (

)

type Data struct {
    data []byte
    index int
    length int
}

func newData (fileData []byte) *Data {
    d := new(Data)
    d.data = fileData
    d.length = len(fileData)
    d.index = -1 //-1 cause it starts with an incriment
    return d
}

func (d *Data) next() (byte, bool){
    d.index++
    if d.index == d.length {
        return ' ', false
    }
    return d.data[d.index], true
}

func (d *Data) peek() (byte, bool) {
    if d.index + 1 == d.length {
        return ' ', false
    }
    return d.data[d.index+1], true
}

func (d *Data) goBack(){
    d.index--
}
