package categories

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *application) routes() http.Handler {
	route := chi.NewRouter()
	route.Use(middleware.Logger)

	route.Get("/", app.getCategoriesListHandler)
	route.Get("/:id", app.getCategoryHandler)
	route.Post("/", app.createCategoryHandler)
	route.Patch("/:id", app.updateCategoryHandler)

	return route
}
