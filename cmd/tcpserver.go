package cmd

import (
	"context"
	"encoding/hex"
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
		log.Fatal("server url does not match HOST:PORT")
	}
	url := strings.Split(serverUrl, ":")
	localUrl := fmt.Sprintf("0.0.0.0:%s", url[1])
	pterm.Info.Println("Starting internal TCP server on", localUrl)
	log.Info("Starting TCP server on ", localUrl)
	// replace ip with 0.0.0.0:port
	TCPServer = NewServer(localUrl)
	TCPServer.sChan = make(chan string, 10000)
}

func (s *Server) ReportIP(conn net.Conn) {
	callbackIP := conn.RemoteAddr().String()

	buf := make([]byte, 4096)
	reqLen, err := conn.Read(buf)
	if err != nil {
		log.Error("Error reading callback buffer", err.Error())
	}

	url := strings.Split(callbackIP, ":")
	msg := fmt.Sprintf("SUCCESS: Remote addr: %s Buf[%d]: %s", url[0], reqLen, hex.EncodeToString(buf[0:reqLen]))
	log.Info(msg)
	pterm.Success.Println(msg)
	if s != nil && s.sChan != nil {
		resMsg := fmt.Sprintf("vulnerable,%s,,", url[0])
		updateCsvRecords(resMsg)
		s.sChan <- resMsg
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
	time.Sleep(10 * time.Second)
	if s == nil || s.listener == nil {
		return
	}
	close(s.quit)
	s.listener.Close()
	spinnerSuccess.Stop()
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
