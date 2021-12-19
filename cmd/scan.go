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
	"strconv"
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
			pterm.Error.Println("server flag error")
			err := cmd.Usage()
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		publicIPAllowed, err := cmd.Flags().GetBool("allow-public-ips")
		if err != nil {
			pterm.Error.Println("allow-public-ip flag error")
			err := cmd.Usage()
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if publicIPAllowed {
			pterm.Warning.Println("Scanning public IPs should be done with care, use at your own risk")
		}
		// TODO: add cancel context
		cidr, err := cmd.Flags().GetString("cidr")
		if err != nil {
			pterm.Error.Println("CIDR flag error")
			err := cmd.Usage()
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if cidr == "" {
			log.Error("CIDR flag missing")
			pterm.Error.Println("CIDR flag missing")
			err := cmd.Usage()
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		CIDRName(cidr)

		ports, err := cmd.Flags().GetString("ports")
		if err != nil {
			pterm.Error.Println("error in ports flag")
			err := cmd.Usage()
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if ports != "top100" && ports != "top10" {
			// check if ports is a single number
			if _, err = strconv.Atoi(ports); err == nil {
			} else {
				// check if ports are ints seperated by colon
				r := strings.Split(ports, ":")
				if len(r) != 2 {
					pterm.Error.Println("error in ports flag")
					cmd.Usage()
					return
				}

				p1, err := strconv.Atoi(r[0])
				if err != nil {
					pterm.Error.Println("error in ports flag")
					cmd.Usage()
					return
				}
				p2, err := strconv.Atoi(r[1])
				if err != nil {
					pterm.Error.Println("error in ports flag")
					cmd.Usage()
					return
				}
				if p2 < p1 {
					pterm.Error.Println("error in ports flag")
					cmd.Usage()
					return
				}
			}
		}

		serverUrl, err := cmd.Flags().GetString("server")
		if err != nil {
			pterm.Error.Println("Error in server flag")
			err := cmd.Usage()
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if serverUrl == "" {
			const port = "5555"
			ipaddrs := GetLocalIP()
			serverUrl = fmt.Sprintf("%s:%s", ipaddrs, port)
		}

		csvPath, err = cmd.Flags().GetString("csv-output")
		if err != nil {
			pterm.Error.Println("Error in csv-output flag")
			err := cmd.Usage()
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		initCSV()

		serverTimeout, err := cmd.Flags().GetInt("timeout")
		if err != nil {
			pterm.Error.Println("error in timeout flag")
			err := cmd.Usage()
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		ctx := context.Background()
		if !disableServer {
			StartServer(ctx, serverUrl, serverTimeout)
		}
		ScanCIDR(ctx, cidr, ports, serverUrl, publicIPAllowed)
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
	scanCmd.Flags().Bool("allow-public-ips", false, "allowing to scan public IPs")
	scanCmd.Flags().String("server", "", "Callback server IP and port (e.g. 192.168.1.100:5555)")
	scanCmd.Flags().String("ports", "top10",
		"Ports to scan. By default scans top 10 ports;"+
			"'top100' will scan the top 100 ports,"+
			"to scan a single insert a port number (e.g. 9000),"+
			"to scan a range of ports, insert range separated by colon (e.g. range 9000:9004) range is limited to max 1024 ports.")
	scanCmd.Flags().String("csv-output", "log4jScanner-results.csv",
		"Set path (inc. filename) to save the CSV file containing the scan results (e.g /tmp/log4jScanner_results.csv). By default will be saved in the running folder.")
	scanCmd.Flags().Int("timeout", 10, "Duration of time to wait before closing the callback server, in secods")
	createPrivateIPBlocks()
}

func ScanCIDR(ctx context.Context, cidr string, portsFlag string, serverUrl string, allowPublicIPs bool) {
	hosts, err := Hosts(cidr, allowPublicIPs)
	//if err is not nil cidr wasn't parse correctly or ip isn't private
	if err != nil {
		pterm.Error.Println("Failed to get hosts, what:", err)
		//an error occurred and program should shut down, close the TCP server
		if LDAPServer != nil {
			LDAPServer.Stop()
		}
		return
	}

	pterm.Info.Printf("Scanning %d addresses in %s\n", len(hosts), cidr)
	// Scan for open ports, if there is an open port, add it to the chan

	// if there are no IPs in the hosts lists, close the TCP server
	if len(hosts) == 0 {
		pterm.Error.Println("No IP addresses in CIDR")
		if LDAPServer != nil {
			LDAPServer.Stop()
		}
		return
	}

	var ports []int
	if portsFlag == "top100" { // top100 will go over to 100 ports
		ports = top100WebPorts
	} else if portsFlag == "top10" {
		ports = top10WebPorts
	} else if p, err := strconv.Atoi(portsFlag); err == nil { // a single port
		ports = append(ports, p)
	} else { // range of ports
		portsRange := strings.Split(portsFlag, ":")
		startPort, _ := strconv.Atoi(portsRange[0])
		endPort, _ := strconv.Atoi(portsRange[1])
		ports = make([]int, endPort-startPort+1)
		for i := range ports {
			ports[i] = startPort + i
			if len(ports) > portRangeSizeLimit {
				pterm.Error.Printfln("port range is limited to %v ports", portRangeSizeLimit)
				if LDAPServer != nil {
					LDAPServer.Stop()
				}
				return
			}
		}
	}

	resChan := make(chan string, 10000)

	var wg sync.WaitGroup
	p, _ := pterm.DefaultProgressbar.WithTotal(len(hosts)).WithTitle("Progress").Start()

	// Loop1: (single go) take all ips, add ports and place in blocking channel, when done close the channel

	// Loop2: (multi go) read ip+port from chan and start a go, once the channel is closed, the loop ends. when done, close the res chan

	// Loop3: (single go) read all results from the res chan, when chan is closed finish

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
	if LDAPServer != nil {
		LDAPServer.Stop()
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

	if LDAPServer != nil && LDAPServer.sChan != nil {
		pterm.Println()
		pterm.NewStyle(pterm.FgGreen).Printfln("Total callbacks: %d", len(LDAPServer.sChan))
		close(LDAPServer.sChan)
		for suc := range LDAPServer.sChan {
			fullSuc := strings.Split(suc, ",")
			msg := fmt.Sprintf("Summary: Callback from %s:%s", fullSuc[1], fullSuc[2])
			pterm.Info.Println(msg)
			log.Info(msg)
		}
	}
}

func ScanPorts(ip, server string, ports []int, resChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Infof("Trying: %s", ip)

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

func Hosts(cidr string, allowPublicIPs bool) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {

		//if public ip scanning isn't allowed Only scan for private IP addresses. If IP is not private, terminate with error.
		if !allowPublicIPs && !isPrivateIP(ip) {
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
