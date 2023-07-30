package main

import (
	"errors"
	"fmt"
	"github.com/go-stack/stack"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
	"net/http"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireBasicAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, plaintextPassword, ok := r.BasicAuth()
		if !ok {
			app.basicAuthenticationRequired(w, r)
			return
		}

		if app.config.basicAuth.username != username {
			app.basicAuthenticationRequired(w, r)
			return
		}

		err := bcrypt.CompareHashAndPassword([]byte(app.config.basicAuth.hashedPassword), []byte(plaintextPassword))
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			app.basicAuthenticationRequired(w, r)
			return
		case err != nil:
			app.serverError(w, r, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Recover(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			defer func() {
				if p := recover(); p != nil {
					err, ok := p.(error)
					if !ok {
						err = fmt.Errorf("%v", p)
					}

					var stackTrace stack.CallStack
					// Get the current stacktrace but trim the runtime
					traces := stack.Trace().TrimRuntime()

					// Format the stack trace removing the clutter from it
					for i := 0; i < len(traces); i++ {
						t := traces[i]
						tFunc := t.Frame().Function

						// Opentelemetry is recovering from the panics on span.End defets and throwing them again
						// we don't want this noise to appear on our logs
						if tFunc == "runtime.gopanic" || tFunc == "go.opentelemetry.io/otel/sdk/trace.(*span).End" {
							continue
						}

						// This call is made before the code reaching our handlers, we don't want to log things that are coming before
						// our own code, just from our handlers and donwards.
						if tFunc == "net/http.HandlerFunc.ServeHTTP" {
							break
						}
						stackTrace = append(stackTrace, t)
					}

					logger.Error(err.Error())

					//logger.WithLevel(zerolog.PanicLevel).
					//	Err(err).
					//	Str("trace.id", trace.SpanFromContext(r.Context()).SpanContext().TraceID().String()).
					//	Str("request-id", GetReqID(r.Context())).
					//	Str("stack", fmt.Sprintf("%+v", stackTrace)).
					//	Msg("panic")

					http.Error(rw, http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(rw, r)
		}
		return http.HandlerFunc(fn)
	}
}
