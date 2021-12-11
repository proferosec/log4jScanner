package main

import (
    log "github.com/sirupsen/logrus"
    "net"
)

func StartServer() {
    // TODO: make port configurable
    const port = ":5555"
    //log.Fatal(http.ListenAndServe(port, r))
    log.Info("Starting server")
    ln, _ := net.Listen("tcp", "192.168.1.59"+port)
    defer ln.Close()
    for {
        log.Debugf("Waiting for connection")
        conn, _ := ln.Accept()
        go ReportIP(conn)
    }

}

func ReportIP(conn net.Conn) {
    log.Infof("SUCCESS: Remote addr: %s", conn.RemoteAddr())
    conn.Close()
}
