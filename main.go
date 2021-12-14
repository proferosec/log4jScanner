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
package main

import (
	log "github.com/sirupsen/logrus"
	"log4jScanner/cmd"
	"log4jScanner/utils"
	"os"
)

var (
	Version   string
	BuildTime string
)

func main() {
	utils.SetVersion(Version)
	//utils.PrintHeader()
	file, err := os.OpenFile("log4jScanner.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	// TODO: fix to enable gosec
	defer file.Close()
	if err != nil {
		log.Error("Failed to log to file")
	}

	utils.InitLogger()
	utils.GetLogger().SetFile(file)
	log.WithFields(log.Fields{"buildTime": BuildTime}).Debugf("Version: ", Version)

	//cmd.SetVersionTemplate("test")
	//cmd.SetHelpFunc()
	cmd.Execute()
}
