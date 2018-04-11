package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/jdel/go-syno"
	"github.com/jdel/gosspks/cfg"
	"github.com/jdel/gosspks/util"
	log "github.com/sirupsen/logrus"
)

var loggerSynology = log.WithFields(log.Fields{
	"module": util.WhereAmI(),
})

// RouteSynology serves packages in Syno compatible JSON format
// This route is meant to serve packages to a Syno
func RouteSynology(w http.ResponseWriter, r *http.Request) {
	var synoPkgs syno.Packages
	var err error
	var scheme = getScheme(r)
	var paramLanguage,
		paramTimezone,
		paramUnique,
		paramArch,
		paramUpdateChannel,
		paramMajor,
		paramMinor,
		paramMicro,
		paramNano,
		paramBuild string

	switch r.Method {
	case http.MethodGet:
		paramLanguage = r.URL.Query().Get("language")
		paramTimezone = r.URL.Query().Get("timezone")
		paramUnique = r.URL.Query().Get("unique")
		paramArch = r.URL.Query().Get("arch")
		paramMajor = r.URL.Query().Get("major")
		paramMinor = r.URL.Query().Get("minor")
		paramMicro = r.URL.Query().Get("micro")
		paramNano = r.URL.Query().Get("nano")
		paramBuild = r.URL.Query().Get("build")
		paramUpdateChannel = r.URL.Query().Get("package_update_channel")
	case http.MethodPost:
		paramLanguage = r.FormValue("language")
		paramTimezone = r.FormValue("timezone")
		paramUnique = r.FormValue("unique")
		paramArch = r.FormValue("arch")
		paramMajor = r.FormValue("major")
		paramMinor = r.FormValue("minor")
		paramMicro = r.FormValue("micro")
		paramNano = r.FormValue("nano")
		paramBuild = r.FormValue("build")
		paramUpdateChannel = r.FormValue("package_update_channel")
	}
	cfg.SynoOptions.Language = paramLanguage

	// Returns bad request if we are missing one of the params
	if paramArch == "" || paramMajor == "" || paramMinor == "" || paramBuild == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if synoPkgs, err = GetPackagesFromCacheOrFileSystem(r, r.URL.Query().Get("language"), false); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		loggerSynology.Warn("Couldn't load packages from cache or filesystem")
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

	var formattedHeader string

	for k, v := range r.Header {
		formattedHeader = fmt.Sprintf("%s \u00A0\u00A0\u00A0 %s: %s \u000A", formattedHeader, k, v[0])
	}

	if cfg.GetDebugPackage() == true {
		debug := syno.NewDebugPackage(fmt.Sprintf(
			// This string cannot be made multiline with backticks
			// because raw string literal is UTF-8 encoded
			"gosspks version: %s \u000A Header: \u000A %s \u000A Arch: %s \u000A Major: %s \u000A Minor: %s \u000A Micro: %s \u000A Nano: %s \u000A Build: %s \u000A Timezone: %s \u000A Language: %s \u000A Unique: %s",
			cfg.Version,
			formattedHeader,
			paramArch,
			paramMajor,
			paramMinor,
			paramMicro,
			paramNano,
			paramBuild,
			paramTimezone,
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
		Keyrings []string       `json:"keyrings,omitempty"`
		Packages *syno.Packages `json:"packages"`
	}{
		Keyrings: cfg.GetGPGKeys(),
		Packages: &synoPkgs,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(jsonOutput); err != nil {
		loggerSynology.Fatal(err)
	}
}
