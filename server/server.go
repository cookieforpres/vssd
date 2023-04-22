package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"vssd/ssd"
)

type RequestMessage struct {
	Action string `json:"action"`
	Size   int    `json:"size"`
	Name   string `json:"name"`
}

type ResponseMessage struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

type Server struct {
	Host    string
	Port    int
	Verbose bool
	SSD     *ssd.SSD
}

func New(host string, port int, name string, size int, verbose bool) *Server {
	return &Server{
		Host:    host,
		Port:    port,
		Verbose: verbose,
		SSD:     ssd.New(name, size),
	}
}

func RespondJSON(message string, erro bool) []byte {
	msg := NewResponseMessage(message, erro)
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return nil
	}

	return data
}

func (ss *Server) Start() error {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ss.Host, ss.Port))
	if err != nil {
		return err
	}

	if ss.Verbose {
		log.Printf("starting tcp server on %s:%d\n", ss.Host, ss.Port)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		if ss.Verbose {
			log.Printf("new connection from %s\n", conn.RemoteAddr())
		}

		go ss.handleConnection(conn)
	}
}

func (ss *Server) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		if ss.Verbose {
			log.Printf("connection from %s closed\n", conn.RemoteAddr())
		}
	}()

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		var jsonData RequestMessage
		err = json.Unmarshal(buf[:n], &jsonData)
		if err != nil {
			log.Println(err)
			return
		}

		switch jsonData.Action {
		case "write":
			buf := make([]byte, jsonData.Size)
			n, err := conn.Read(buf)
			if err != nil {
				log.Println(err)
				return
			}

			err = ss.SSD.Write(ssd.NewNode(jsonData.Name), buf[:n])
			if err != nil {
				log.Println(err)
				return
			}

			_, err = conn.Write(RespondJSON("ok", false))
			if err != nil {
				log.Println(err)
				return
			}
		case "read":
			data, err := ss.SSD.Read(jsonData.Name)
			if err != nil {
				if err == ssd.ErrNoNode {
					_, err = conn.Write(RespondJSON("no node", true))
					if err != nil {
						log.Println(err)
						return
					}

					continue
				}

				log.Println(err)
				return
			}

			_, err = conn.Write(data)
			if err != nil {
				log.Println(err)
				return
			}

		case "delete":
			err := ss.SSD.Delete(jsonData.Name)
			if err != nil {
				log.Println(err)
				return
			}

			_, err = conn.Write(RespondJSON("ok", false))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func NewResponseMessage(message string, err bool) *ResponseMessage {
	return &ResponseMessage{
		Message: message,
		Error:   err,
	}
}

func GB(size int) int {
	return size * 1024 * 1024 * 1024
}

func MB(size int) int {
	return size * 1024 * 1024
}

func KB(size int) int {
	return size * 1024
}

func B(size int) int {
	return size
}
