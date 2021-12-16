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
	"context"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan all IPs in the given CIDR",
	Long: `Scan each IP for open ports. By default will scan 10 top ports.
For example: log4jScanner scan --cidr "192.168.0.1/24`,
	Run: func(cmd *cobra.Command, args []string) {
		nocolors, err := cmd.Flags().GetBool("nocolor")
		if nocolors {
			pterm.DisableColor()
		}

		disableServer, err := cmd.Flags().GetBool("noserver")
		if err != nil {
			log.Error("server flag error")
			cmd.Usage()
			return
		}
		// TODO: add cancel context
		cidr, err := cmd.Flags().GetString("cidr")
		if err != nil {
			log.Error("CIDR flag error")
			cmd.Usage()
			return
		}
		if cidr == "" {
			log.Error("CIDR flag missing")
			pterm.Error.Println("CIDR flag missing")
			cmd.Usage()
			return
		}
		CIDRName(cidr)

		ports, err := cmd.Flags().GetString("ports")
		if err != nil || (ports != "top100" && ports != "slow" && ports != "top10") {
			fmt.Println("error in ports flag")
			cmd.Usage()
			return
		}

		serverUrl, err := cmd.Flags().GetString("server")
		if err != nil {
			fmt.Println("Error in server flag")
			cmd.Usage()
			return
		}
		if serverUrl == "" {
			const port = "5555"
			ipaddrs := GetLocalIP()
			serverUrl = fmt.Sprintf("%s:%s", ipaddrs, port)
		}

		csvPath, err = cmd.Flags().GetString("csv-output")
		if err != nil {
			fmt.Println("Error in csv-output flag")
			cmd.Usage()
			return
		}
		initCSV()

		ctx := context.Background()
		if !disableServer {
			StartServer(ctx, serverUrl)
		}
		ScanCIDR(ctx, cidr, ports, serverUrl)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	scanCmd.Flags().String("cidr", "", "IP subnet to scan in CIDR notation (e.g. 192.168.1.0/24)")
	scanCmd.Flags().Bool("noserver", false, "Do not use the internal TCP server, this overrides the server flag if present")
	scanCmd.Flags().Bool("nocolor", false, "remove colors from output")
	scanCmd.Flags().String("server", "", "Callback server IP and port (e.g. 192.168.1.100:5555)")
	scanCmd.Flags().String("ports", "top10",
		"Ports to scan. By default scans top 10 ports; 'top100' will scan the top 100 ports, 'slow' will scan all possible ports")
	scanCmd.Flags().String("csv-output", "log4jScanner-results.csv",
		"Set path (inc. filename) to save the CSV file containing the scan results (e.g /tmp/log4jScanner_results.csv). By default will be saved in the running folder.")
	createPrivateIPBlocks()
}

func ScanCIDR(ctx context.Context, cidr string, portsFlag string, serverUrl string) {
	hosts, err := Hosts(cidr)
	//if err is not nil cidr wasn't parse correctly or ip isn't private
	if err != nil {
		pterm.Error.Println("Failed to get hosts, what:", err)
		//an error occurred and program should shut down, close the TCP server
		if TCPServer != nil {
			TCPServer.Stop()
		}
		return
	}

	pterm.Info.Printf("Scanning %d addresses in %s\n", len(hosts), cidr)
	// Scan for open ports, if there is an open port, add it to the chan

	// if there are no IPs in the hosts lists, close the TCP server
	if len(hosts) == 0 {
		pterm.Error.Println("No IP addresses in CIDR")
		if TCPServer != nil {
			TCPServer.Stop()
		}
		return
	}

	var ports []int
	if portsFlag == "slow" {
		pterm.Warning.Println("Slow flag is currently disabled, defaulting to top10")
		ports = top10WebPorts
		//ports = make([]int, endPortSlow-startPortSlow+1)
		//for i := range ports {
		//	ports[i] = startPortSlow + i
		//}
	} else if portsFlag == "top100" { // top100 will go over to 100 ports
		ports = top100WebPorts
	} else { // Fast scan - will go over the ports from the top 10 ports list.
		ports = top10WebPorts
	}

	resChan := make(chan string, 10000)

	var wg sync.WaitGroup
	p, _ := pterm.DefaultProgressbar.WithTotal(len(hosts)).WithTitle("Progress").Start()
	const maxGoroutines = 100
	cnt := 0
	for _, i := range hosts {
		cnt += 1
		// TODO: replace ports flag with an ENUM
		if cnt > maxGoroutines {
			wg.Wait()
			cnt = 0
		}
		wg.Add(1)
		p.Increment()
		// TODO: replace with go
		ScanPorts(i, serverUrl, ports, resChan, &wg)
	}
	wg.Wait()
	if TCPServer != nil {
		TCPServer.Stop()
	}
	PrintResults(resChan)
}

func PrintResults(resChan chan string) {
	close(resChan)
	pterm.Println()
	pterm.NewStyle(pterm.FgGreen).Printfln("Total requests: %d", len(resChan))
	for res := range resChan {
		fullRes := strings.Split(res, ",")
		msg := fmt.Sprintf("Summary: %s:%s ==> %s", fullRes[1], fullRes[2], fullRes[3])
		pterm.Info.Println(msg)
		log.Info(msg)
	}

	if TCPServer != nil && TCPServer.sChan != nil {
		pterm.Println()
		pterm.NewStyle(pterm.FgGreen).Printfln("Total callbacks: %d", len(TCPServer.sChan))
		close(TCPServer.sChan)
		for suc := range TCPServer.sChan {
			fullSuc := strings.Split(suc, ",")
			msg := fmt.Sprintf("Summary: Callback from %s", fullSuc[1])
			pterm.Info.Println(msg)
			log.Info(msg)
		}
	}
}

func ScanPorts(ip, server string, ports []int, resChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Infof("Trying: %s", ip)

	// Slow scan will go over all ports from 1 to 65535
	wgPorts := sync.WaitGroup{}
	for _, port := range ports {
		targetHttps := fmt.Sprintf("http://%s:%v", ip, port)
		targetHttp := fmt.Sprintf("https://%s:%v", ip, port)
		wgPorts.Add(2)
		go ScanIP(targetHttp, server, &wgPorts, resChan)
		go ScanIP(targetHttps, server, &wgPorts, resChan)
	}
	wgPorts.Wait()

}

func Hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		// Only scan for private IP addresses. If IP is not private, skip.
		if !isPrivateIP(ip) {
			badIPStatus := ip.String() + " IP address is not private"
			pterm.Error.Println(badIPStatus)
			log.Fatal(badIPStatus)
		}
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

func isPrivateIP(ip net.IP) bool {
	//ip := net.ParseIP(ipS)

	for _, block := range privateIPs {
		if block.Contains(ip) {
			return true
		}
	}
	return false

}

func createPrivateIPBlocks() {
	for _, cidr := range privateIPBlocks {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			log.Error("parse error on %q: %v", cidr, err)
		}
		privateIPs = append(privateIPs, block)
	}
}

//func isPortHttps(port int) bool {
//	for _, p := range commonHttpsPorts {
//		if port == p {
//			return true
//		}
//	}
//	return false
//}
