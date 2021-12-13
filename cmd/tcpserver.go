package cmd

import (
    "context"
    "fmt"
    "github.com/pterm/pterm"
    log "github.com/sirupsen/logrus"
    "io"
    "net"
    "strings"
    "sync"
    "time"
)

var TCPServer *Server

func StartServer(ctx context.Context, server_url string) {
    // TODO: make port configurable
    if server_url != "" && !strings.Contains(server_url, ":") {
        pterm.Error.Println("server url does not match HOST:PORT")
        return
    }
    if server_url == "" {
        const port = "5555"
        ipaddrs := GetLocalIP()
        server_url = fmt.Sprintf("%s:%s", ipaddrs, port)
    }
    pterm.Info.Println("Starting internal TCP server on", server_url)
    log.Info("Starting TCP server on", server_url)
    TCPServer = NewServer(server_url)
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
    if s == nil || s.listener == nil {
        return
    }
    close(s.quit)
    s.listener.Close()
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
