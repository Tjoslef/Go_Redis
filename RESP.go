package main
import (
	"fmt"
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

const(
	STRING = '+'
	ERROR = '-'
	ARRAY = '*'
	INTEGER = ':'
	BULK = '$'
)
type Value struct{
	typ string
	num int
	bulk string
	array []Value
}
type Resp struct {
		reader *bufio.Reader
}
func NewResp(rd io.Reader) *Resp (
	return &Resp{reader : bufio.NewReader(rd)}
)
func (r * Resp) readLine() (line[]  byte, n int, err error) {
	for {
		b, err := r.reader.Readbyte()
		if err != nil {
			return nil,0,err
		}
		n += 1
		line = append(line,b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
	}




	}

		return line[:len(line)-2], n, nil
}
func (r * Resp) readInteger() (x int, n int, err error){
	
		line,n,err := r.reader.readLine()
		if err != nil{
			return 0,0,err
		}
		i64,err := strconv.ParseInt(string(line),10,64)
		if err != nil {
			return 0,n,err
	}
		return int(i64),n, nil
	}
func (r* Resp) Read() (Value,error) {
	_type, err := r.reader.Readbyte()

	if err != nil {
		return Value{}, err
	}
	switch _type {
	case ARRAY:
		return r.readArray()
	case Bulk:
		return r.readBulk()
	default:
		fmt.Println("invalid type")
		return Value{},nil

	}
}
func (r*Resp) readArray() (Value, error){
	v := Value{}
	v.typ = "array"
	size,_,err := r.readInteger()
	if err != nil{
		return v,err
	}
	v.array = make([]Value, 0)
	for i := 0;i < size;i++{
		val,err := r.Read()
		if err != nil {
			return v,err
		}
		v.array = append(v.array,val)

	}
	return v,nil
}
func (r*Resp) readBulk () (Value,error) {
	v := Value{}
	v.typ = "bulk"
	_,size,err := r.readLine()
	if err != nil {
		return v, err
	}
	val := make([]byte,size)
	r.reader.Read(val)
	v.bulk = string(val)
	return v,nil

}

	
