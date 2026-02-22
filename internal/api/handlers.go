package api

import (
	"net/http"
	"notas-app-go/internal/data"
)

type Application struct {
	Notas data.NotaModel
}

func (app *Application) GetStatus(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, http.StatusOK, nil)
}
