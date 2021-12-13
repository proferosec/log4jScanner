package cmd

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
)

var TCPServer *Server

func StartServer(ctx context.Context, serverUrl string) {
	// TODO: make port configurable
	if serverUrl != "" && !strings.Contains(serverUrl, ":") {
		pterm.Error.Println("server url does not match HOST:PORT")
		return
	}
	pterm.Info.Println("Starting internal TCP server on", serverUrl)
	log.Info("Starting TCP server on ", serverUrl)
	TCPServer = NewServer(serverUrl)
	TCPServer.sChan = make(chan string, 10000)
}

func (s *Server) ReportIP(conn net.Conn) {
	msg := fmt.Sprintf("SUCCESS: Remote addr: %s", conn.RemoteAddr())
	log.Info(msg)
	pterm.Success.Println(msg)
	if s != nil && s.sChan != nil {
		s.sChan <- msg
	}
	conn.Close()
}

type Server struct {
	listener net.Listener
	quit     chan interface{}
	wg       sync.WaitGroup
	sChan    chan string
}

func NewServer(addr string) *Server {
	s := &Server{
		quit: make(chan interface{}),
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		pterm.Error.Println(err)
		log.Fatal(err)
	}
	s.listener = l
	s.wg.Add(1)
	go s.serve()
	return s
}

func (s *Server) Stop() {
	spinnerSuccess, _ := pterm.DefaultSpinner.Start("Stopping TCP server")
	time.Sleep(2 * time.Second)
	if s == nil || s.listener == nil {
		return
	}
	close(s.quit)
	s.listener.Close()
	spinnerSuccess.Success()
	s.wg.Wait()
}

func (s *Server) serve() {
	defer s.wg.Done()

	for {
		if s == nil || s.listener == nil {
			return
		}
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				log.Println("accept error", err)
			}
		} else {
			s.ReportIP(conn)
		}
	}
}
