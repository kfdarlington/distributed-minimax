package web

import (
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"github.com/kristian-d/distributed-minimax/engine/leader/pools"
	"net/http"
)

type route struct {
	method string
	endpoint string
	handler http.HandlerFunc
}

type handler struct {
	pools *pools.Pool
	logger *logger.Logger
}

type routes []route

func handlerize(fn func (http.ResponseWriter, *http.Request, interface{})) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		fn(w, r, nil)
	}
}

func NewRouter(pools *pools.Pool, logger *logger.Logger) http.Handler {
	var h = &handler{
		pools: pools,
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
			"/followers",
			handlerize(h.postFollowers),
		},
		route{
			"GET",
			"/followers",
			handlerize(h.getFollowers),
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
