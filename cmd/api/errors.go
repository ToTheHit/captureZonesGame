package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"main/internal/response"
	"main/internal/validator"

	"golang.org/x/exp/slog"
)

func (app *application) reportError(err error) {
	trace := string(debug.Stack())
	app.logger.Error(err.Error(), slog.String("trace", trace))

	if app.config.notifications.email != "" {
		data := app.newEmailData()
		data["Timestamp"] = time.Now().UTC().Format(time.RFC3339)
		data["Error"] = err.Error()
		data["Trace"] = trace

		err := app.mailer.Send(app.config.notifications.email, data, "error-notification.tmpl")
		if err != nil {
			trace := string(debug.Stack())
			app.logger.Error(err.Error(), slog.String("trace", trace))
		}
	}
}

func (app *application) errorMessage(w http.ResponseWriter, r *http.Request, status int, message string, headers http.Header) {
	message = strings.ToUpper(message[:1]) + message[1:]

	err := response.JSONWithHeaders(w, status, map[string]string{"Error": message}, headers)
	if err != nil {
		app.reportError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.reportError(err)

	message := "The server encountered a problem and could not process your request"
	app.errorMessage(w, r, http.StatusInternalServerError, message, nil)
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	app.errorMessage(w, r, http.StatusNotFound, message, nil)
}

func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	app.errorMessage(w, r, http.StatusMethodNotAllowed, message, nil)
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.errorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
}

func (app *application) failedValidation(w http.ResponseWriter, r *http.Request, v validator.Validator) {
	err := response.JSON(w, http.StatusUnprocessableEntity, v)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) basicAuthenticationRequired(w http.ResponseWriter, r *http.Request) {
	headers := make(http.Header)
	headers.Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	message := "You must be authenticated to access this resource"
	app.errorMessage(w, r, http.StatusUnauthorized, message, headers)
}
