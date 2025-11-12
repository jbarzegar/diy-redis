package main

import (
	"diy-redis/internal/server"
	"io"
	"log"
)

func main() {
	srv, err := server.NewServer()
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Conn.Close()

	for {
		_, err = srv.Read()
		if err != nil {
			if err == io.EOF {
				continue
			}
			log.Fatal(err)
		}
		command, err := srv.GetCommand()
		if err != nil {
			srv.Conn.Write([]byte("+ParsingError\r\n"))
			continue
		}

		srv.HandleArgs(command)
		continue
	}
}
