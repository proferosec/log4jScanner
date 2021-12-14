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
	url := strings.Split(serverUrl, ":")
	localUrl := fmt.Sprintf("0.0.0.0:%s", url[1])
	pterm.Info.Println("Starting internal TCP server on", localUrl)
	log.Info("Starting TCP server on ", localUrl)
	// replace ip with 0.0.0.0:port
	TCPServer = NewServer(localUrl)
	TCPServer.sChan = make(chan string, 10000)
	TCPServer.csvHandler = csvHandler
}

func csvHandler(tpe, ip, port string) {
	rec := []string{tpe, ip, port}
	fmt.Println(rec)
}

func (s *Server) ReportIP(conn net.Conn) {
	callbackIP := conn.RemoteAddr().String()
	msg := fmt.Sprintf("SUCCESS: Remote addr: %s", callbackIP)
	url := strings.Split(callbackIP, ":")
	log.Info(msg)
	pterm.Success.Println(msg)
	if s != nil && s.sChan != nil {
		if s.csvHandler != nil {
			s.csvHandler("callback", url[0], url[1])
		}
		s.sChan <- fmt.Sprintf("callback,%s,%s,", url[0], url[1])
	}
	conn.Close()
}

type Server struct {
	listener   net.Listener
	quit       chan interface{}
	wg         sync.WaitGroup
	sChan      chan string
	csvHandler func(tpe, ip, port string)
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
