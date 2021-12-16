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
	"os"

	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
	"time"

	"log4jScanner/utils"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

var LogPath string
var DebugFlag bool
var logPathFlag string
var CIDR string
var logTime string

const logDateFormat = "2006-01-02_150405"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "log4jScanner",
	Short: "Root command",
	Long: `log4jScanner tool will scan a subnet for web servers and will try to send the JNDI exploit to each one. 
			For every response it receives, it will log the sender IP so we can get a list of the vulnerable servers.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(utils.PrintHeader)
	//initLog need to be run after header is been printed for output order
	cobra.OnInitialize(initLog)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.log4j_scanner.yaml)")

	rootCmd.PersistentFlags().BoolVar(&DebugFlag, "debug", false, "set log level to debug")
	rootCmd.PersistentFlags().StringVar(&logPathFlag, "log-output", "", "Set name and path to save the log file (e.g  /tmp/log4jScanner.log). By default will be saved in the running folder")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("server", "s", false, "Run callback server")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	logTime = time.Now().Format(logDateFormat)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".log4j_scanner" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".log4j_scanner")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

//initLog check if --debug flag was set
//and will set the logger to debug level
//otherwise logger level will be info by default
func initLog() {
	if DebugFlag {
		utils.Logger.SetLevel(log.DebugLevel)
	}

	// log name includes time & CIDR flag
	if logPathFlag == "" && LogPath == "" {
		LogPath = fmt.Sprintf("log4jScanner-%s-%s.log", CIDR, logTime)
	} else if logPathFlag == "" {
		LogP := fmt.Sprintf("log4jScanner-%s-%s.log", CIDR, logTime)
		err := os.Rename(LogPath, LogP)
		if err != nil {
			log.Error(err)
		}
		LogPath = LogP
	} else {
		lSuffix := filepath.Ext(logPathFlag)
		LogP := fmt.Sprintf("%s-%s-%s%s", strings.TrimSuffix(logPathFlag, lSuffix), CIDR, logTime, lSuffix)
		if LogPath != "" {
			err := os.Rename(LogPath, LogP)
			if err != nil {
				log.Error(err)
			}
		}
		LogPath = LogP
	}
	file, err := os.OpenFile(LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		pterm.Warning.Println("failed to change log file location (using running folder), what:", err)
		file, err = os.OpenFile("log4jScanner.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			log.Fatal("Failed to log to file, what:", err)
		}
	}
	utils.GetLogger().SetFile(file)
}

// Format CIDR to add to log name
func CIDRName(cidr string) {
	cidr = strings.ReplaceAll(cidr, ".", "_")
	cidr = strings.ReplaceAll(cidr, "/", "__")
	CIDR = cidr

	initLog()
}
