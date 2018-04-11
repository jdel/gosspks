package v1

import (
	"encoding/json"
	"net/http"

	"github.com/jdel/gosspks/util"
	log "github.com/sirupsen/logrus"
)

var loggerModels = log.WithFields(log.Fields{
	"module": util.WhereAmI(),
})

// RouteModels displays synology models
func RouteModels(w http.ResponseWriter, r *http.Request) {
	synoModels := GetModelsFromCacheOrWeb(false)
	if r.URL.Query().Get("name") != "" {
		synoModels = synoModels.FilterByName(r.URL.Query().Get("name"))
	}

	if synoModels == nil {
		w.WriteHeader(http.StatusNotFound)
		loggerModels.Warn("No models available")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(synoModels); err != nil {
		loggerModels.Fatal(err)
	}
}
