package bootstrap

import (
	"fmt"
	"log"
	"net/http"
)

func (app *Application) Serve() error {
	serverAddr := fmt.Sprintf(":%d", app.Port)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: app.Server,
	}

	log.Printf("Starting server on port %d", app.Port)

	return server.ListenAndServe()
}
