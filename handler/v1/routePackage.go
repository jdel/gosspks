package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"jdel.org/go-syno"
	"jdel.org/gosspks/cfg"
	"jdel.org/gosspks/util"
	log "github.com/sirupsen/logrus"
)

var loggerPackage = log.WithFields(log.Fields{
	"module": util.WhereAmI(),
})

// RoutePackage serves a single packages in JSON format
// This route is not meant to serve packages to a syno
func RoutePackage(w http.ResponseWriter, r *http.Request) {
	cfg.SynoOptions.Language = r.URL.Query().Get("language")

	if synoPkg, err := syno.NewPackage(mux.Vars(r)["synoPackageName"]); err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(synoPkg); err != nil {
			loggerPackage.Fatal(err)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
