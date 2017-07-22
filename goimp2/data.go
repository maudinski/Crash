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
// arbitrarily chose @ to mean EOF. Will be fine within strings tho
func (d *Data) next() string {
    d.index++
    if d.index == d.length {
        return "@"
    }
    return string(d.data[d.index])
}

func (d *Data) peek() string {
    if d.index + 1 == d.length {
        return "@"
    }
    return string(d.data[d.index+1])
}

func (d *Data) goBack(){
    d.index--
}
