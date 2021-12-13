package cmd

import (
	"context"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
)

var TCPServer *Server

func StartServer(ctx context.Context, serverUrl string, successChan chan string) {
    // TODO: make port configurable
    if serverUrl != "" && !strings.Contains(serverUrl, ":") {
        pterm.Error.Println("server url does not match HOST:PORT")
        return
    }
    pterm.Info.Println("Starting internal TCP server on", serverUrl)
    log.Info("Starting TCP server on ", serverUrl)
    TCPServer = NewServer(serverUrl)
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
        pterm.Error.Println(err)
        log.Error(err)
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
    s.wg.Wait()
    spinnerSuccess.Success()
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
    conn.SetDeadline(time.Now().Add(1000 * time.Millisecond))
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
