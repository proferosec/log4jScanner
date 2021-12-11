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
    joonix "github.com/joonix/log"
    log "github.com/sirupsen/logrus"
    "log4j_scanner/cmd"
    "os"
    "strings"
)

var (
    Version   string
    BuildTime string
)

// TODO: add version to the version command

func setupLog(logFormat, logLevel string) {
    switch strings.ToLower(logFormat) {
    case "text":
        {
            log.SetFormatter(&log.TextFormatter{}) // normal output
        }
    case "json":
        {
            log.SetFormatter(&log.JSONFormatter{}) // simple json output
        }
    case "fluentd":
        {
            log.SetFormatter(joonix.NewFormatter()) //Fluentd compatible
        }
    default:
        log.SetFormatter(joonix.NewFormatter()) //Fluentd compatible
    }

    log.SetOutput(os.Stdout)
    switch strings.ToLower(logLevel) {
    case "debug":
        log.SetLevel(log.DebugLevel)
    case "warning":
        log.SetLevel(log.WarnLevel)
    default:
        log.SetLevel(log.InfoLevel)
    }
}

// TODO: log to file
// TODO: add header/pterm

func main() {
    setupLog("text", "debug")
    log.WithFields(log.Fields{"buildTime": BuildTime}).Info("Version: ", Version)

    go cmd.Execute()

    StartServer()
}
