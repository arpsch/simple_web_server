package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	api_http "ws/api/http"
)

func readServerPort() string {

	port := flag.String("port", "8080", "a port for the server to listen, or defaults to :8080")
	flag.Parse()

	return ":" + *port
}

func setupServer() (*httprouter.Router, error) {

	appHandler := api_http.NewAppHandlers()
	routes := appHandler.SetupRoutes()

	return routes, nil
}

func main() {
	router, err := setupServer()
	if err != nil {
		log.Printf(" failed to set up routes, exiting")
		return
	}

	port := readServerPort()
	log.Printf("Server is listening on %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
