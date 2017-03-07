package handlers

import (
	"net/http"

	"github.com/wobeproject/logger"
)

//Route ..
type route struct {
	name        string
	pattern     string
	handlerFunc func(http.ResponseWriter, *http.Request)
	handlers    []func(http.Handler) http.Handler
}

var routes = []route{
	route{
		"index",
		"/",
		IndexHandler,
		[]func(http.Handler) http.Handler{
			ValidationHandler, RecoverHandler,
		},
	},
}

//addRoutes ...
func addRoutes(routes ...route) {
	for _, r := range routes {
		l.Debug("Routes", map[string]interface{}{
			"Adding routes": r.name,
		})
		var handler http.Handler
		handler = http.HandlerFunc(r.handlerFunc)
		for i := range r.handlers {
			handler = r.handlers[i](handler)
		}
		http.Handle(r.pattern, handler)
	}
}

//InitHandlers ...
func InitHandlers() {
	l = logger.GetInstance()
	addRoutes(routes...)

}
