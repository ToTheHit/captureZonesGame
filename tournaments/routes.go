package tournaments

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *application) routes() http.Handler {
	route := chi.NewRouter()
	route.Use(middleware.Logger)
	//tournamentsRoutes, err := buildTournamentRoutes()
	//
	//route.Mount("/DDD", tournamentRoutes)

	return route

}
