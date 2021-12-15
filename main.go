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
)

var (
    Version   string
    BuildTime string
)

func main() {
    utils.SetVersion(Version, BuildTime)
    //utils.PrintHeader()

    utils.InitLogger()
    defer utils.Logger.Close()
    log.WithFields(log.Fields{"buildTime": BuildTime}).Debugf("Version: ", Version)

    //cmd.SetVersionTemplate("test")
    //cmd.SetHelpFunc()
    cmd.Execute()
}
