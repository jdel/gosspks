package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"github.com/jdel/go-syno"
	"github.com/jdel/gosspks/cfg"
	"github.com/jdel/gosspks/util"
	log "github.com/sirupsen/logrus"
)

var loggerSynofficial = log.WithFields(log.Fields{
	"module": util.WhereAmI(),
})

// RouteSynofficial serves packages in Syno compatible JSON format
// This route is meant to serve packages to a Syno emulating the official
// Synology package repository route.
// To use this route, add an entry pointing to this server for
// pkgautoupdate.synology.com in your Syno's host file
func RouteSynofficial(w http.ResponseWriter, r *http.Request) {
	var synoPkgs syno.Packages
	var err error
	var scheme = getScheme(r)
	var paramLanguage,
		paramUnique,
		paramArch,
		paramUpdateChannel,
		paramMajor,
		paramMinor,
		paramMicro,
		paramNano,
		paramBuild string

	// Get mux vars from URL
	cfg.SynoOptions.Language = r.URL.Query().Get("language")
	paramLanguage = mux.Vars(r)["synoLanguage"]
	paramUnique = mux.Vars(r)["synoUnique"]
	paramArch = mux.Vars(r)["synoArch"]
	paramMajor = mux.Vars(r)["synoMajor"]
	paramMinor = mux.Vars(r)["synoMinor"]
	paramMicro = mux.Vars(r)["synoMicro"]
	paramNano = mux.Vars(r)["synoNano"]
	paramBuild = mux.Vars(r)["synoBuild"]
	paramUpdateChannel = mux.Vars(r)["synoChannel"]

	if synoPkgs, err = GetPackagesFromCacheOrFileSystem(r, r.URL.Query().Get("language"), false); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		loggerSynofficial.Warn("Couldn't load packages from cache or filesystem")
	}

	// Filter and sort packages
	synoPkgs = synoPkgs.
		FilterByArch(paramArch).
		FilterByFirmware(paramMajor + "." + paramMinor + "-" + paramBuild).
		OnlyShowLastVersion()

	if paramUpdateChannel != "beta" {
		synoPkgs = synoPkgs.FilterOutBeta()
	}

	sort.Sort(synoPkgs)

	if cfg.GetDebugPackage() == true {
		debug := syno.NewDebugPackage(fmt.Sprintf(
			"Header: %s, Arch: %s, Major: %s, Minor: %s, Micro: %s, Nano: %s, Build: %s, Language: %s, Unique: %s",
			r.Header,
			paramArch,
			paramMajor,
			paramMinor,
			paramMicro,
			paramNano,
			paramBuild,
			paramLanguage,
			paramUnique))
		debug.Thumbnail = append(debug.Thumbnail, fmt.Sprintf(
			"%s://%s/%s/debug.png",
			scheme, cfg.GetStaticPrefix(), r.Host))
		synoPkgs = append(synoPkgs, debug)
	}

	// Wraps synoPkgs in a packages[] json array
	// No support for synology build < 5004 (DSM 5.1 Update 2)
	jsonOutput := struct {
		Keyrings *[]string      `json:"keyrings,omitempty"`
		Packages *syno.Packages `json:"packages"`
	}{
		Keyrings: nil,
		Packages: &synoPkgs,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(jsonOutput); err != nil {
		loggerSynofficial.Fatal(err)
	}
}
