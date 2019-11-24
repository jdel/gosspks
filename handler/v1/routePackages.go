package v1

import (
	"encoding/json"
	"net/http"
	"sort"

	"jdel.org/go-syno"
	"jdel.org/gosspks/cfg"
	"jdel.org/gosspks/util"
	log "github.com/sirupsen/logrus"
)

var loggerPackages = log.WithFields(log.Fields{
	"module": util.WhereAmI(),
})

// RoutePackages serves all available packages in JSON format
// This route is not meant to serve packages to a syno
func RoutePackages(w http.ResponseWriter, r *http.Request) {
	var synoPkgs syno.Packages
	var err error
	var paramArch,
		paramName,
		paramLatest string

	cfg.SynoOptions.Language = r.URL.Query().Get("language")
	paramArch = r.URL.Query().Get("arch")
	paramName = r.URL.Query().Get("name")
	paramLatest = r.URL.Query().Get("latest")

	if synoPkgs, err = GetPackagesFromCacheOrFileSystem(r, r.URL.Query().Get("language"), false); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		loggerPackages.Warn("Couldn't load packages from cache or filesystem")
	}

	if len(synoPkgs) == 0 {
		w.WriteHeader(http.StatusNotFound)
		loggerPackages.Warn("No packages available in ", cfg.GetPackagesDir())
		return
	}

	// Filter and sort packages
	if paramArch != "" {
		synoPkgs = synoPkgs.FilterByArch(paramArch)
	}
	if paramLatest == "true" {
		synoPkgs = synoPkgs.OnlyShowLastVersion()
	}
	if paramName != "" {
		synoPkgs = synoPkgs.SearchByName(paramName)
	}

	sort.Sort(synoPkgs)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(synoPkgs); err != nil {
		loggerPackages.Fatal(err)
	}
}
