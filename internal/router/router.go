package router

import (
	"github.com/EpicStep/avito-autumn-2021-intern-task/internal/jsonutil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

// New returns new chi Mux.
func New() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		jsonutil.MarshalResponse(w, http.StatusNotFound, jsonutil.NewError(3, "API method not found"))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		jsonutil.MarshalResponse(w, http.StatusMethodNotAllowed, jsonutil.NewError(3, "HTTP method not allowed"))
	})

	return r
}
