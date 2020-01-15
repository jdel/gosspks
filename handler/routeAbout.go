package handler // import jdel.org/gosspks/handler

import (
	"time"
	"encoding/json"
	"net/http"

	"jdel.org/gosspks/cfg"
	"jdel.org/gosspks/util"
	log "github.com/sirupsen/logrus"
)

var loggerAbout = log.WithFields(log.Fields{
	"module": util.WhereAmI(),
})

// RouteAbout displays some information about
// the running instance of gosspks
func RouteAbout(w http.ResponseWriter, r *http.Request) {
	about := struct {
		Name       string `json:"name"`
		Version    string `json:"version"`
		Maintainer string `json:"maintainer"`
		License    string `json:"license"`
		Year       int    `json:"year"`
	}{
		"gosspks",
		cfg.Version,
		"jdel",
		"GNU GPL v3",
		time.Now().Year(),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(about); err != nil {
		loggerAbout.Fatal(err)
	}
}
