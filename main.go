package main

import (
	respparser "diy-redis/internal/resp"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Command struct {
	Name string
}

func CmdOK() []byte {
	return []byte("+OK\r\n")
}

func printRawCmd(b []byte) {
	fmt.Printf("%q\n", strings.Trim(string(b), "\x00"))
}

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
		fmt.Printf("%q\n", str)
		reader := respparser.NewRespParser(strings.NewReader(str))
		result, err := reader.Read()
		if err != nil {
			conn.Write([]byte("+ParsingError\r\n"))
			continue
		}

		s := fmt.Sprintf("+%v(%v)\r\n", result.Kind, result.Args)
		conn.Write([]byte(s))

		// conn.Write(CmdOK())
		continue
	}
}
