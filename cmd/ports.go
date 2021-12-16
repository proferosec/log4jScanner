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
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// portsCmd represents the ports command
var portsCmd = &cobra.Command{
	Use:   "ports",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		PrintPorts()
	},
}

func PrintPorts() {
	pterm.DefaultCenter.Println("Top 10 Ports")
	for _, p := range top10WebPorts {
		pterm.Println(p)
	}
	pterm.DefaultCenter.Println("Top 100 Ports")
	for _, p := range top100WebPorts {
		pterm.Println(p)
	}
}

func init() {
	rootCmd.AddCommand(portsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// portsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// portsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
