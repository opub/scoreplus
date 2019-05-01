package web

import (
	"net/http"

	"github.com/go-chi/render"
)

// StatusResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type StatusResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

//Render pre-processes response
func (e *StatusResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

//ErrBadRequest renderer for invalid syntax
func ErrBadRequest(err error) render.Renderer {
	return &StatusResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Bad Request",
		ErrorText:      err.Error(),
	}
}

//ErrServerError renderer internal application errors
func ErrServerError(err error) render.Renderer {
	return &StatusResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal Server Error",
		ErrorText:      err.Error(),
	}
}

//ErrUnauthorized renderer for 401
var ErrUnauthorized = &StatusResponse{HTTPStatusCode: 1, StatusText: "Not Authorized"}

//ErrNotFound renderer for 404
var ErrNotFound = &StatusResponse{HTTPStatusCode: 404, StatusText: "Not Found"}

//StatusOK renderer for success
var StatusOK = &StatusResponse{HTTPStatusCode: 200, StatusText: "Success"}
