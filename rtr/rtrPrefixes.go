package rtr // import jdel.org/gosspks/rtr

import (
	"net/http"

	"jdel.org/gosspks/cfg"
	"jdel.org/gosspks/handler"
)

type pathPrefix struct {
	Name        string
	Methods     string
	Pattern     func() string
	HandlerFunc http.HandlerFunc
}

type pathPrefixes []pathPrefix

// This variable is responsible for configuring all prefixes
// prefix funcs are in handlers/prefixes.go
var appPrefixes = pathPrefixes{
	pathPrefix{
		Name:        "Static",
		Methods:     "GET",
		Pattern:     func() string { return cfg.GetStaticPrefix() },
		HandlerFunc: handler.PrefixStatic,
	},
	pathPrefix{
		Name:        "Download",
		Methods:     "GET",
		Pattern:     func() string { return cfg.GetDownloadPrefix() },
		HandlerFunc: handler.PrefixDownload,
	},
}
