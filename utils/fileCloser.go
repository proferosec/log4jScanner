package utils

import (
    log "github.com/sirupsen/logrus"
    "os"
)

//FileCloser is a closing a file and checking the error
//this function should be used with defer keyword
//using "defer f.close" create a gosec error
func FileCloser(f *os.File) {
    if err := f.Close(); err != nil {
        log.Fatal(err)
    }
}
