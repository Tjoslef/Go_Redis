package main

const {
	STRING = '+'
	ERROR = '-'
	ARRAY = '*'
	INTEGER = ':'
	BULK = '$'
	}
type Value struct{
	typ string
	num int
	bulk string
	array []Value
}
import (
	"fmt"
	"bufio"
	"string"
	"os"
	"strconv"


)
type Resp struct {
	reader *bufio.Reader
}
func NewResp(rd io.Reader) *Resp (
	return &Resp{reader : bufio.NewReader(rd)}
)
func (r * Resp) readLine (line[]  byte, n int, err error) {
	for {
		b, err := r.reader.Readbyte()
		if err != nil {
			return nil,0,err
		}
		n += 1
		line = append(line,b]
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
	}




	}

		return line[:len(line)-2], n, nil
}
	
