package main

import (
	respparser "diy-redis/internal/resp"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	// start a network connection
	//
	ln, err := net.Listen("tcp", ":6379")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening on port :6379")
	conn, err := ln.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		// make a byte array of the max TCP packet size
		buf := make([]byte, 1024)

		_, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				continue
			}
			log.Fatal(err)
		}
		str := strings.Trim(string(buf), "\x00")
		reader := respparser.NewRespParser(strings.NewReader(str))
		command, err := reader.Read()
		if err != nil {
			conn.Write([]byte("+ParsingError\r\n"))
			continue
		}

		if command.Kind == respparser.DatatypeArray {
			headCmd := command.Args[0]
			handleCommand(conn, headCmd)
		}
		continue
	}
}

func handleCommand(con net.Conn, c respparser.Command, args ...[]respparser.Command) {
	fmt.Println(c.Value, c)
	switch c.Value {
	case "PING":
		con.Write([]byte("+\"PONG\"\r\n"))
	default:
		con.Write([]byte("-\"UNSUPPORTED\"\r\n"))
	}
}
