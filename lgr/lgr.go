package lgr

import (
	"jdel.org/gosspks/util"
	log "github.com/sirupsen/logrus"
)

// LogrusWriter provides a io.Writer implementation
// to logrus
type LogrusWriter struct {
	Module   string
	Route    string
	Template string
}

func (w *LogrusWriter) Write(b []byte) (int, error) {
	// remove traling \n
	n := len(b)
	if n > 0 && b[n-1] == '\n' {
		b = b[:n-1]
	}

	log.WithFields(log.Fields{
		"module":   util.WhereAmI(),
		"route":    w.Route,
		"template": w.Template,
	}).Info(string(b))
	return n, nil
}
