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

//addRoutes ...
func addRoutes(routes ...route) {
	for _, r := range routes {
		var handler http.Handler
		handler = http.HandlerFunc(r.handlerFunc)
		for _, h := range r.handlers {
			handler = h(handler)
		}
		http.Handle(r.pattern, handler)
	}
}

//InitHandlers ...
func InitHandlers() {
	l = logger.GetInstance()
	addRoutes(
		route{
			"index",
			"/",
			InputHandler,
			[]func(http.Handler) http.Handler{
				RecoverHandler, ValidationHandler,
			},
		},
	)
}
