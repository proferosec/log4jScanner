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
)

// TODO: update all descriptions

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ScanCIDR()
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ScanCIDR() {
	hosts, _ := Hosts("192.168.1.58/28")
	ipsChan := make(chan string, 1024)
	ipPortChan := make(chan string, 256)
	//doneChan := make(chan string)

	const concurrentMax = 100

	// Scan for open ports, if there is an open port, add it to the chan
	for _, ip := range hosts {
		ipsChan <- ip
	}

	for i := range ipsChan {
		ScanPorts(i, ipPortChan, false)
		if len(ipsChan) == 0 {
			close(ipsChan)
			close(ipPortChan)
		}
	}
}

func ScanPorts(ip string, ipPortChan chan string, slow bool) {
	log.Debugf("Trying: %s", ip)
	var ports []int

	// Slow scan will go over all ports from 1 to 63556
	if slow {
		log.Debugln("Slow scan")
		ports = make([]int, endPortSlow-startPortSlow+1)
		for i := range ports {
			ports[i] = startPortSlow + i
		}
		// Fast scan - will go over the ports from the input ports list.
	} else {
		ports = topWebPorts
	}

	localIP := GetLocalIP()
	log.Debugf("Local IP: %s", localIP)
	go ScanIP(ipPortChan, localIP+":5555")

	for _, port := range ports {
		target := fmt.Sprintf("http://%s:%v", ip, port)
		ipPortChan <- target
	}

}

func Hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	// remove network address and broadcast address
	lenIPs := len(ips)
	switch {
	case lenIPs < 2:
		return ips, nil

	default:
		return ips[1 : len(ips)-1], nil
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
