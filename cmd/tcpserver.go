package cmd

import (
	"context"
	"fmt"
	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"sync"
	"time"
)

var TCPServer *Server

func StartServer(ctx context.Context) {
	// TODO: make port configurable
	const port = "5555"
	//log.Fatal(http.ListenAndServe(port, r))
	log.Info("Starting server")
	//ln, _ := net.Listen("tcp", "192.168.1.59"+port)
	//defer ln.Close()
	//for {
	//    if StopServer {
	//        break
	//    }
	//    conn, _ := ln.Accept()
	//    go ReportIP(conn)
	//}
	ipaddrs := GetLocalIP()
	TCPServer = NewServer(fmt.Sprintf("%s:%s", ipaddrs, port))
}

func ReportIP(conn net.Conn) {
	log.Infof("SUCCESS: Remote addr: %s", conn.RemoteAddr())
	pterm.Success.Println("Remote address: ", conn.RemoteAddr())
	conn.Close()
}

type Server struct {
	listener net.Listener
	quit     chan interface{}
	wg       sync.WaitGroup
}

func NewServer(addr string) *Server {
	s := &Server{
		quit: make(chan interface{}),
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	s.listener = l
	s.wg.Add(1)
	go s.serve()
	return s
}

func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
}

func (s *Server) serve() {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				log.Println("accept error", err)
			}
		} else {
			s.wg.Add(1)
			go func() {
				s.handleConection(conn)
				s.wg.Done()
			}()
		}
	}
}

func (s *Server) handleConection(conn net.Conn) {
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(200 * time.Millisecond))
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			log.Println("read error", err)
			return
		}
		if n == 0 {
			return
		}
		ReportIP(conn)
	}
}
