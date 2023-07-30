package questions

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *application) routes() http.Handler {
	route := chi.NewRouter()
	route.Use(middleware.Logger)

	route.Get("/", app.getQuestionsListHandler)
	route.Get("/:id", app.getQuestionHandler)
	route.Post("/", app.createQuestionHandler)
	route.Patch("/:id", app.updateQuestionHandler)

	return route
}
