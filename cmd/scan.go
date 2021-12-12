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
	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"log4jScanner/utils"
	"net"
	"sync"
)

// TODO: update all descriptions
// TODO: update the header for all usages
// TODO: add commandline flags

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
		utils.PrintHeader()
		enableServer, err := cmd.Flags().GetBool("server")
		if err != nil {
			log.Error("server flag error")
		}
		// TODO: add cancel context
		cidr, err := cmd.Flags().GetString("cidr")
		if err != nil || cidr == "" {
			log.Error("CIDR flag missing")
			cmd.Usage()
			return
		}

		slow, err := cmd.Flags().GetBool("slow")
		if err != nil {
			log.Error("slow flag error")
		}

		ctx := context.Background()
		ServerStartOnFlag(ctx, enableServer)
		ScanCIDR(ctx, cidr, slow)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	scanCmd.PersistentFlags().String("cidr", "", "IP subnet to scan in CIDR notation (e.g. 192.168.1.0/24)")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	scanCmd.Flags().BoolP("server", "s", false, "Use internal TCP server")

	scanCmd.Flags().Bool("slow", false, "Slow scan will scan all ports between 1-65535")

	createPrivateIPBlocks()
}

func ServerStartOnFlag(ctx context.Context, enable bool) {
	if enable {
		pterm.Info.Println("Starting internal TCP server")
		StartServer(ctx)
	}
}

func ScanCIDR(ctx context.Context, cidr string, slow bool) {
	hosts, _ := Hosts(cidr)
	ipsChan := make(chan string, 1024)
	ipPortChan := make(chan string, 256)
	//doneChan := make(chan string)

	pterm.Info.Printf("Scanning %d addresses in %s\n", len(hosts), cidr)
	// Scan for open ports, if there is an open port, add it to the chan
	for _, ip := range hosts {
		// Only scan for private IP addresses. If IP is not private, skip.
		if !isPrivateIP(ip) {
			log.Errorf("%s IP adress is not private", ip)
			continue
		}
		ipsChan <- ip
	}

	if len(ipsChan) == 0 {
		close(ipsChan)
		if TCPServer != nil {
			TCPServer.Stop()
		}
	}

	server := GetLocalIP() + ":5555"

	var wg sync.WaitGroup
	p, _ := pterm.DefaultProgressbar.WithTotal(len(ipsChan)).WithTitle("Progress").Start()
	for i := range ipsChan {
		wg.Add(1)
		go ScanPorts(i, server, ipPortChan, slow, p, &wg)
		if len(ipsChan) == 0 {
			close(ipsChan)
		}
	}
	wg.Wait()
	if TCPServer != nil {
		TCPServer.Stop()
	}
}

func ScanPorts(ip, server string, ipPortChan chan string, slow bool, p *pterm.ProgressbarPrinter, wg *sync.WaitGroup) {
	var ports []int

	log.Infof("Trying: %s", ip)

	// Slow scan will go over all ports from 1 to 65535
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

	for _, port := range ports {
		target := fmt.Sprintf("http://%s:%v", ip, port)
		ipPortChan <- target
	}
	go ScanIP(ipPortChan, server)
	p.Increment()
	wg.Done()
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

func isPrivateIP(ipS string) bool {
	ip := net.ParseIP(ipS)

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
