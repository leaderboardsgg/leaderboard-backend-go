package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/speedrun-website/speedrun-rest/data"
)

func loadPlatformRoutes(router *mux.Router) {
	platformRouter := router.PathPrefix("/platforms").Subrouter()

	// Index
	platformRouter.
		Path("/").
		HandlerFunc(getPlatforms).
		Name("PlatformIndex").
		Methods("GET")

	// Show
	platformRouter.
		Path("/{name}").
		HandlerFunc(getPlatformByName).
		Name("PlatformShow").
		Methods("GET")

	// Create
	platformRouter.
		Path("/").
		HandlerFunc(authenticatedHandler(createPlatform)).
		Name("PlatformCreate").
		Methods("POST")
}

func getPlatforms(w http.ResponseWriter, r *http.Request) {
	platforms := data.GetPlatforms()
	writeJsonResponse(w, platforms)
}

func getPlatformByName(w http.ResponseWriter, r *http.Request) {
	platformName := mux.Vars(r)["name"]
	platform, err := data.GetPlatformByName(platformName)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	writeJsonResponse(w, platform)
}

func createPlatform(w http.ResponseWriter, r *http.Request) {
	var platform data.Platform
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&platform)

	data.CreatePlatform(platform)

	writeJsonResponse(w, platform)
}
