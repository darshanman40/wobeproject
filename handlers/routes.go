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
}

//AddRoutes ...
func addRoutes(routes ...route) {
	for _, r := range routes {
		var handler http.Handler
		handler = http.HandlerFunc(r.handlerFunc) //handlers.InputHandler)
		handler = recoverHandler(handler)

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
			inputHandler,
		},
	)
}
