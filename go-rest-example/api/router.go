package api

import (
	"github.com/gorilla/mux"
)

func LoadRouter() *mux.Router {
	router := mux.NewRouter()

	loadPlatformRoutes(router)
	// loadGameRoutes(router)

	return router
}
