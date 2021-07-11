package main

import (
	"fmt"
	"net/http"

	"github.com/speedrun-website/speedrun-rest/api"
)

func main() {
	router := api.LoadRouter()

	fmt.Println("listening on 8080")
	http.ListenAndServe(":8080", router)
}
