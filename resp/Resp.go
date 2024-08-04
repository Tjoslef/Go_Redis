package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	Typ   string
	Str   string
	num   int
	Bulk  string
	Array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}

func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.Typ = "array"

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	v.Array = make([]Value, 0)
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		v.Array = append(v.Array, val)
	}

	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}

	v.Typ = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.Bulk = string(bulk)

	r.readLine()

	return v, nil
}
func (v Value) Marshal() []byte {
	switch v.Typ {
	case "array":

		return v.writerArray()
	case "bulk":
		return v.writerBulk()
	case "string":

		return v.writerString()
	case "null":
		return v.writerNull()
	case "error":

		return v.writerError()
	default:
		return []byte{}
	}
}
func (v Value) writerArray() []byte {
	size := len(v.Array)
	var byteArray = make([]byte, size)
	byteArray = append(byteArray, ARRAY)
	byteArray = append(byteArray, strconv.Itoa(size)...)
	byteArray = append(byteArray, '\r', '\n')
	for i := 0; i < size; i++ {
		byteArray = append(byteArray, v.Array[i].Marshal()...)

	}
	return byteArray
}
func (v Value) writerBulk() []byte {
	size := len(v.Bulk)
	var byteBulk []byte
	byteBulk = append(byteBulk, BULK)
	byteBulk = append(byteBulk, strconv.Itoa(size)...)
	byteBulk = append(byteBulk, '\r', '\n')
	byteBulk = append(byteBulk, v.Bulk...)
	byteBulk = append(byteBulk, '\r', '\n')
	return byteBulk
}
func (v Value) writerString() []byte {
	var byteString []byte
	byteString = append(byteString, STRING)
	byteString = append(byteString, v.Str...)
	byteString = append(byteString, '\r', '\n')
	return byteString

}
func (v Value) writerNull() []byte {
	return []byte("$-1\r\n")
}
func (v Value) writerError() []byte {
	var byteError []byte
	byteError = append(byteError, ERROR)
	byteError = append(byteError, v.Str...)
	byteError = append(byteError, '\r', '\n')
	return byteError
}

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) Write(v Value) error {
	var bytes = v.Marshal()
	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
