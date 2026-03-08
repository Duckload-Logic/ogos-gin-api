package bootstrap

import (
	"fmt"
	"net/http"
)

func (app *Application) Serve() error {
	serverAddr := fmt.Sprintf(":%d", app.Port)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: app.Server,
	}

	return server.ListenAndServe()
}
