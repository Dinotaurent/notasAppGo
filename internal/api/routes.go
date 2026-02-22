package api

import "net/http"

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	// Rutas
	mux.HandleFunc("/status", app.GetStatus)

	return app.Logger(mux)
}
