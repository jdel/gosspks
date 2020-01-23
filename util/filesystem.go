package util // import jdel.org/gosspks/util

import (
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"
)

var loggerFS = log.WithFields(log.Fields{
	"module": WhereAmI(),
})

// FileInfoBySizeDesc implementats a sort
// by size on os.FileInfo
type FileInfoBySizeDesc []os.FileInfo

func (s FileInfoBySizeDesc) Len() int {
	return len(s)
}
func (s FileInfoBySizeDesc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s FileInfoBySizeDesc) Less(i, j int) bool {
	return s[i].Size() > s[j].Size()
}

// WhereAmI returns the current source file name 
func WhereAmI() string {
	_, filename, _, _ := runtime.Caller(1)
	return filename
}

// CreateDir Creates a directory if it doesn't exist
func CreateDir(dir string) error {
	var err error
	if !FileExists(dir) {
		loggerFS.Debug("Creating directory: ", dir)
		err = os.MkdirAll(dir, 0755)
	}
	return err
}

// FileExists returns true if the dir exists on filesystem
func FileExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}
