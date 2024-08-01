package main
import (
	"fmt"
	"io"
	"net"
	"os"

)
func main() {

	fmt.Println("Start listing on port :6379")
	// creating a server
	l,err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return 
	}
	conn,err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		buff := make([]byte,1024)
		_,err := conn.Read(buff)
		if err != nil {
			if err == io.EOF{
				break
			}
		fmt.Println("error reading from client", err.Error())	
		os.Exit(1)
		}
		conn.Write([]byte("+OK\r\n"))



	}
}

