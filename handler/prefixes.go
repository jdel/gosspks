package handler // import jdel.org/gosspks/handler

import (
	"net/http"
	"strings"

	"jdel.org/gosspks/cfg"
)

// PrefixStatic serves anything under /static/
// (or whatever gosspks.routes.static is set to)
func PrefixStatic(w http.ResponseWriter, r *http.Request) {
	// Translates prefix from URL to local FS (/static/ -> cache/)
	localPath := strings.Replace(r.URL.Path, cfg.GetStaticPrefix(), cfg.GetCacheDir(), 1)
	http.ServeFile(w, r, localPath)
}

// PrefixDownload serves anything under /download/
// (or whatever gosspks.routes.download is set to)
func PrefixDownload(w http.ResponseWriter, r *http.Request) {
	// Translates prefix from URL to local FS (/download/ -> packages/)
	localPath := strings.Replace(r.URL.Path, cfg.GetDownloadPrefix(), cfg.GetPackagesDir(), 1)
	http.ServeFile(w, r, localPath)
}
