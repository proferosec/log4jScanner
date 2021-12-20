package cmd

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	ldap "github.com/vjeantet/ldapserver"

	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
)

var LDAPServer *Server

func StartServer(ctx context.Context, serverUrl string, serverTimeout int) {
	pterm.Info.Println("Server URL: " + serverUrl)
	log.Info("Server URL: " + serverUrl)
	listenUrl, err := url.Parse("//" + serverUrl)
	if err != nil {
		pterm.Error.Println("Failed to parse server url")
		log.Fatal("Failed to parse server url")
	}
	// replace ip with 0.0.0.0:port
	listenUrl.Host = "0.0.0.0:" + listenUrl.Port()

	pterm.Info.Println("Starting internal LDAP server on", listenUrl.Host)
	log.Info("Starting LDAP server on ", listenUrl.Host)
	LDAPServer = NewServer()
	LDAPServer.sChan = make(chan string, 10000)
	LDAPServer.timeout = time.Duration(serverTimeout) * time.Second

	go LDAPServer.server.ListenAndServe(listenUrl.Host)
}

func (s *Server) ReportIP(vulnerableServiceLocation string) {
	vulnUrl, err := url.Parse("//" + vulnerableServiceLocation)
	if err != nil {
		pterm.Error.Println("Failed to parse vulnerable url: " + vulnerableServiceLocation)
		log.Fatal("Failed to parse server url" + vulnerableServiceLocation)
	}
	msg := fmt.Sprintf("SUCCESS: Remote addr: %s:%s", vulnUrl.Hostname(), vulnUrl.Port())
	log.Info(msg)
	pterm.Success.Println(msg)
	if s != nil && s.sChan != nil {
		resMsg := fmt.Sprintf("vulnerable,%s,%s,", vulnUrl.Hostname(), vulnUrl.Port())
		updateCsvRecords(resMsg)
		s.sChan <- resMsg
	}
}

type Server struct {
	server  *ldap.Server
	sChan   chan string
	timeout time.Duration
}

func (s *Server) handleBind(w ldap.ResponseWriter, m *ldap.Message) {
	res := ldap.NewBindResponse(ldap.LDAPResultSuccess)
	w.Write(res)
	return
}

func (s *Server) handleSearch(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetSearchRequest()

	//pterm.Info.Println("Got LDAP search request: " + r.BaseObject())
	log.Info("Got LDAP search request: " + r.BaseObject())

	vulnerableLocation := strings.ReplaceAll(string(r.BaseObject()), "_", ":")

	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
	w.Write(res)

	s.ReportIP(vulnerableLocation)

	return
}

func NewServer() *Server {
	s := &Server{
		server: ldap.NewServer(),
	}

	ldap.Logger = log.StandardLogger()

	routes := ldap.NewRouteMux()
	routes.Bind(s.handleBind)
	routes.Search(s.handleSearch)

	s.server.Handle(routes)

	return s
}

func (s *Server) Stop() {
	spinnerSuccess, _ := pterm.DefaultSpinner.Start("Stopping LDAP server")
	timeout := LDAPServer.timeout
	time.Sleep(timeout)
	s.server.Stop()
	err := spinnerSuccess.Stop()
	if err != nil {
		log.Fatal(err)
	}
}
