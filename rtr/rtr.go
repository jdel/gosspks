package rtr // import jdel.org/gosspks/rtr

import (
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"jdel.org/gosspks/lgr"
	"jdel.org/gosspks/util"
	log "github.com/sirupsen/logrus"
)

var loggerRouter = log.WithFields(log.Fields{
	"module": util.WhereAmI(),
})

// NewRouter is the mux router that handle main routes
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	// Create parh routes from appRoutes in rtr/routes.go
	for _, route := range appRoutes {
		router.
			Methods(strings.Split(
				strings.Replace(
					route.Methods, " ", "", -1), ",")...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handlers.LoggingHandler(&lgr.LogrusWriter{
				Route:    route.Name,
				Template: route.Pattern}, route.HandlerFunc))
	}
	// Create path prefix routes from appPrefixes in rtr/routes.go
	for _, prefix := range appPrefixes {
		router.
			Methods(strings.Split(
				strings.Replace(
					prefix.Methods, " ", "", -1), ",")...).
			PathPrefix(prefix.Pattern()).
			Name(prefix.Name).
			Handler(handlers.LoggingHandler(&lgr.LogrusWriter{
				Route:    prefix.Name,
				Template: prefix.Pattern()}, prefix.HandlerFunc))
	}
	// Prints routes in debug log
	router.Walk(printRoutes)
	return router
}

// Prints all available routes for debugging purpose
func printRoutes(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	template, err := route.GetPathTemplate()
	if err != nil {
		return err
	}
	methods, err := route.GetMethods()
	if err != nil {
		return err
	}
	loggerRouter.WithFields(log.Fields{
		"name":     route.GetName(),
		"methods":  strings.Join(methods, ","),
		"template": template,
	}).Debug("Route details")
	return nil
}
