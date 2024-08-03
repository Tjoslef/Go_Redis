package main

import (
	"fmt"
	"net"
	"strings"
	"tjoslef/skola/Redis/resp"
)

func main() {
	fmt.Println("Listening on port :6379")

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		rd := resp.NewResp(conn)
		value, err := rd.Read()
		if err != nil {
			fmt.Println(err)
			return
		}
		if value.Typ != "array" {
			fmt.Println("wrong typ expected array")
		}
		if len(value.Array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}
		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		writer := resp.NewWriter(conn)
		handler, ok := resp.Handler[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(resp.Value{Typ: "string", Str: "OK"})
			continue
		}
		result := handler(args)
		writer.Write(result)

	}

}
