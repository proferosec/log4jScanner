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
    "encoding/csv"
    "fmt"
    "os"
    "strings"

    "github.com/pterm/pterm"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "run a local TCPServer server",
    Long:  "",
    Run: func(cmd *cobra.Command, args []string) {
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
        initCSV()
        ServerStart(serverUrl)
        //pterm.Error.Println("Functionality disabled")
    },
}

func init() {
    //rootCmd.AddCommand(serverCmd)
    serverCmd.Flags().String("server", "", "Callback server IP and port (e.g. 192.168.1.100:5555)")
}

func ServerStart(serverUrl string) {
    StartServer(nil, serverUrl)
    pterm.Info.Println("Press ctr-l-c to exit")
    for {
    }
    pterm.DefaultHeader.WithFullWidth().Println("Results")
    csvRecords := [][]string{
        {"type", "ip", "port", "status_code"},
    }
    PrintServerResults(csvRecords)

}

func PrintServerResults(csvRecords [][]string) {
    close(TCPServer.sChan)
    for suc := range TCPServer.sChan {
        csvSuc := strings.Split(suc, ",")
        msg := fmt.Sprintf("Summary: Callback from %s:%s", csvSuc[1], csvSuc[2])
        pterm.Info.Println(msg)
        log.Info(msg)
        csvRecords = append(csvRecords, csvSuc)
    }
    f, err := os.Create("log4jScanner-results.csv")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    w := csv.NewWriter(f)
    defer w.Flush()

    for _, record := range csvRecords {
        if err := w.Write(record); err != nil {
            log.Fatal(err)
        }
    }

}
