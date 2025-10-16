package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
)

func main() {
	conn, err := net.Dial("tcp", "irc.libera.chat:6667")
	handle_error(err)
	defer conn.Close()

	fmt.Fprintf(conn, "NICK newtoIrc")
	fmt.Fprintf(conn, "USER newtoIrc 0 * : Name")

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	for {

		data, err := tp.ReadLine()
		handle_error(err)
		log.Println(data)
	}
}

func handle_error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
