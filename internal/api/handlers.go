package api

import (
	"encoding/json"
	"math"
	"net/http"
	"notas-app-go/internal/data"
	"strconv"
)

type Application struct {
	Notas data.NotaModel
}

func (app *Application) GetStatus(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, http.StatusOK, nil)
}

func (app *Application) GetNotas(w http.ResponseWriter, r *http.Request) {
	// Obtenemos parámetros: ?page=1&size=10
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	size, _ := strconv.ParseInt(r.URL.Query().Get("size"), 10, 64)

	// Valores por defecto
	if page < 1 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	notas, total, err := app.Notas.GetAll(page, size)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	// Calculamos número de páginas
	var pages int64 = 1
	if total > 0 {
		pages = int64(math.Ceil(float64(total) / float64(size)))
	}

	resp := map[string]any{
		"data":  notas,
		"total": total,
		"page":  page,
		"pages": pages,
	}
	app.writeJSON(w, http.StatusOK, resp)
}

func (app *Application) CreateNota(w http.ResponseWriter, r *http.Request) {
	var input data.Nota

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	if err := app.Notas.Insert(input); err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	app.writeJSON(w, http.StatusCreated, nil)
}

func (app *Application) GetNotaByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	nota, err := app.Notas.GetByID(id)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			app.errorJSON(w, err, http.StatusNotFound)
		}
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	app.writeJSON(w, http.StatusOK, nota)
}

func (app *Application) UpdateNotaByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input data.Nota
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	if err := app.Notas.Update(id, input); err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusAccepted, nil)
}

func (app *Application) DeleteNotaByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := app.Notas.Delete(id); err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusAccepted, nil)
}
