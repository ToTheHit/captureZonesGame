package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"main/categories"
	"main/questions"
	"main/tournaments"
	"net/http"
	"time"
)

func (app *application) routes() http.Handler {
	//mux := httprouter.New()
	//
	//mux.NotFound = http.HandlerFunc(app.notFound)
	//mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)
	//
	//mux.HandlerFunc("GET", "/status", app.status)
	//
	//mux.Handler("GET", "/basic-auth-protected", app.requireBasicAuthentication(http.HandlerFunc(app.protected)))

	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	//r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	//r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(10 * time.Second))

	r.Use(ArticleCtx)

	r.Get("/status", app.status)
	r.With(app.requireBasicAuthentication).Get("/basic-auth-protected", app.protected)

	r.Route("/api", func(r chi.Router) {
		categoriesRoutes, err := categories.BuildRoutes()
		if err != nil {
			app.logger.Error(err.Error())
		}
		questionsRoutes, err := questions.BuildRoutes()
		if err != nil {
			app.logger.Error(err.Error())
		}
		tournamentsRoutes, err := tournaments.BuildRoutes()
		if err != nil {
			app.logger.Error(err.Error())
		}

		r.Use(Recover(app.logger))
		r.Mount("/categories", categoriesRoutes)
		r.Mount("/questions", questionsRoutes)
		r.Mount("/tournaments", tournamentsRoutes)
	})

	return app.recoverPanic(r)
}

func ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Вытаскивает параметр из урла. Пример: GET /{articleID}/status
		articleID := chi.URLParam(r, "articleID")
		fmt.Println("articleID:", articleID)

		// Вытаскивает параметр из query. Пример: GET /status?articleID=test
		query := r.URL.Query().Get("articleID")
		fmt.Println("query:", query)

		ctx := context.WithValue(r.Context(), "article", articleID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
