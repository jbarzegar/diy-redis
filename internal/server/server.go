package server

import (
	respparser "diy-redis/internal/resp"
	"diy-redis/internal/store"
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	store     *store.Store
	cmdReader *respparser.RespParser
	Conn      net.Conn
	buf       []byte
}

// func (*s Server)
func NewServer() (*Server, error) {
	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		// log.Fatal(err)
		return &Server{}, err
	}
	// make a byte array of the max TCP packet size
	buf := make([]byte, 1024)
	s := store.NewStore()

	fmt.Println("Listening on port :6379")
	conn, err := ln.Accept()
	if err != nil {
		log.Fatal(err)
	}
	return &Server{
		store: s,
		buf:   buf,
		Conn:  conn,
	}, nil
}

func (s *Server) Read() (int, error) {
	return s.Conn.Read(s.buf)
}

func (s *Server) GetCommand() (respparser.Command, error) {
	str := strings.Trim(string(s.buf), "\x00")
	reader := respparser.NewRespParser(strings.NewReader(str))
	command, err := reader.Read()
	if err != nil {
		return respparser.Command{}, err
	}

	return command, nil
}

func (s *Server) HandleArgs(command respparser.Command) {
	headCmd := command.Args[0]
	args := command.Args[1:]
	response := s.HandleCommand(headCmd, args)
	s.Conn.Write(response)
}

func printIntResponse(n int) []byte {
	return fmt.Appendf([]byte{}, "+(integer) %v\r\n", n)
}

func (s *Server) HandleCommand(c respparser.Command, args []respparser.Command) []byte {
	switch strings.ToUpper(c.Value) {
	case "PING":
		return []byte("+\"PONG\"\r\n")
	case "SET":
		key := args[0].Value
		val := args[1].Value
		s.store.Set(key, val)
		return []byte("+\"OK\"\r\n")
	case "GET":
		item, _ := s.store.Get(args[0].Value)
		return fmt.Appendf([]byte{}, "+\"%v\"\r\n", item.Value)
	case "EXISTS":
		vals := []string{c.Value}
		for _, a := range args {
			vals = append(vals, a.Value)
		}
		count, _ := s.store.Exists(vals...)
		return printIntResponse(count)
	case "DEL":
		keys := []string{c.Value}
		for _, a := range args {
			keys = append(keys, a.Value)
		}
		count, _ := s.store.Del(keys...)
		return printIntResponse(count)
	default:
		return []byte("-\"UNSUPPORTED\"\r\n")
	}
}
