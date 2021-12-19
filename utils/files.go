package utils

import (
    log "github.com/sirupsen/logrus"
    "io/ioutil"
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

//RenameFile tries to rename the file,
//on success return (nil,nil)
//on failure return (nil,err) if can't read file content, (data,nil) if was able to read file content
func RenameFile(oldPath string, newPath string) ([]byte, error) {
    err := os.Rename(oldPath,newPath)
    //if rename was successful do nothing
    if err == nil {
        return nil, nil
    }
    GetLogger().File.Close()
    data, err := ioutil.ReadFile(oldPath)
    if err != nil {
        return nil, err
    }
    err = os.Remove(oldPath)
    if err != nil {
        return nil, err
    }
    return data, nil
}