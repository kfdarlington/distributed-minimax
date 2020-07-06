package web

import (
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"github.com/kristian-d/distributed-minimax/engine/leader"
	"net/http"
)

type route struct {
	method string
	endpoint string
	handler http.HandlerFunc
}

type handler struct {
	engine *leader.Leader
	logger *logger.Logger
}

type routes []route

func handlerize(fn func (http.ResponseWriter, *http.Request, interface{})) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		fn(w, r, nil)
	}
}

func NewRouter(engine *leader.Leader, logger *logger.Logger) http.Handler {
	var h = &handler{
		engine: engine,
		logger: logger,
	}

	var myRoutes = routes{
		route{
			"GET",
			"/",
			handlerize(h.index),
		},
		route{
			"GET",
			"/ping",
			handlerize(h.ping),
		},
		route{
			"POST",
			"/start",
			handlerize(h.start),
		},
		route{
			"POST",
			"/move",
			handlerize(h.move),
		},
		route{
			"POST",
			"/end",
			handlerize(h.end),
		},
	}

	router := mux.NewRouter().StrictSlash(false)
	for _, route := range myRoutes {
		router.
			Methods(route.method).
			Path(route.endpoint).
			Handler(route.handler)
	}

	return router
}
