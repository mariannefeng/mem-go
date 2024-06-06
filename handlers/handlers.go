package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
)

type Route interface {
	Path() string
	Handler() func(http.ResponseWriter, *http.Request) error
	Methods() []string
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}

type HTTPError struct {
	statusCode int
	err        error

	Message string `json:"message"`
}

func NewHTTPError(statusCode int, message string, err error) *HTTPError {
	return &HTTPError{
		statusCode: statusCode,
		err:        err,
		Message:    message,
	}
}

func (e HTTPError) Error() string {
	return e.Message
}

func (e HTTPError) Unwrap() error {
	return e.err
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func NewMuxRouter(routes []Route) *mux.Router {
	r := mux.NewRouter()

	for _, route := range routes {
		r.HandleFunc(route.Path(), errorHandler(route.Handler())).Methods(route.Methods()...)
	}
	return r
}

func errorHandler(h HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")

				errResp, _ := json.Marshal(&HTTPError{
					Message: "panik",
				})
				w.Write(errResp)
			}
		}()

		err := h(w, r)

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			var httpError *HTTPError

			switch err := err.(type) {
			case HTTPError:
				w.WriteHeader(err.statusCode)
				httpError = &err
			default:
				w.WriteHeader(http.StatusInternalServerError)
				httpError = &HTTPError{
					Message: "something went terribly wrong",
				}
			}

			errResp, _ := json.Marshal(httpError)
			w.Write(errResp)
		}
	}
}
