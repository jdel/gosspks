package rtr

import (
	"net/http"

	"github.com/jdel/gosspks/handler"
	"github.com/jdel/gosspks/handler/v1"
)

type route struct {
	Name        string
	Methods     string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

// This variable is responsible for configuring all routes
var appRoutes = routes{
	route{
		Name:        "About",
		Methods:     "GET",
		Pattern:     "/about",
		HandlerFunc: handler.RouteAbout,
	},
	route{
		Name:        "Package",
		Methods:     "GET",
		Pattern:     "/getList/v0/{synoMajor}/{synoMinor}/{synoMicro}/{synoBuild}/{synoNano}/{synoArch}/{synoChannel}/{synoUnique}/{synoLanguage}",
		HandlerFunc: v1.RouteSynofficial,
	},
	route{
		Name:        "Models",
		Methods:     "GET",
		Pattern:     "/v1/models",
		HandlerFunc: v1.RouteModels,
	},
	route{
		Name:        "Package",
		Methods:     "GET",
		Pattern:     "/v1/packages/{synoPackageName}",
		HandlerFunc: v1.RoutePackage,
	},
	route{
		Name:        "Packages",
		Methods:     "GET",
		Pattern:     "/v1/packages",
		HandlerFunc: v1.RoutePackages,
	},
	route{
		Name:        "SynoPackages",
		Methods:     "POST",
		Pattern:     "/",
		HandlerFunc: v1.RouteSynology,
	},
	route{
		Name:        "SynoPackages",
		Methods:     "GET",
		Pattern:     "/",
		HandlerFunc: v1.RouteSynology,
	},
}
