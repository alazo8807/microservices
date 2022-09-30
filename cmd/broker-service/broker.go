package main

import (
	"broker/internal/routes"
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type Config struct {
}

func main() {
	log.Printf("Starting broker service on port %s", webPort)

	// deifne http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: routes.InitRoutes(),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
