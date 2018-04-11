package rtr

import (
	"net/http"

	"github.com/jdel/gosspks/cfg"
	"github.com/jdel/gosspks/handler"
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
