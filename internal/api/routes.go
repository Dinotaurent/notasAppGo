package api

import "net/http"

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	// Rutas
	mux.HandleFunc("GET /status", app.GetStatus)
	mux.HandleFunc("GET /notas", app.GetNotas)
	mux.HandleFunc("GET /notas/{id}", app.GetNotaByID)
	mux.HandleFunc("POST /notas", app.CreateNota)
	mux.HandleFunc("PUT /notas/{id}", app.UpdateNotaByID)
	mux.HandleFunc("DELETE /notas/{id}", app.DeleteNotaByID)

	return app.Logger(mux)
}
