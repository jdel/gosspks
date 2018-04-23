package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jdel/gosspks/cfg"
	"github.com/jdel/gosspks/util"
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
		Year       uint   `json:"year"`
	}{
		"gosspks",
		cfg.Version,
		"jdel",
		"GNU GPL v3",
		2018,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(about); err != nil {
		loggerAbout.Fatal(err)
	}
}
