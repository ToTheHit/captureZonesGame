package main

import (
	"net/http"

	"main/internal/response"
)

func (app *application) status(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	article := ctx.Value("article").(string)

	data := map[string]string{
		"Status": article,
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) protected(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a protected handler"))
}
