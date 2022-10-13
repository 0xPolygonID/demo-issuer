package rest

import (
	"github.com/go-chi/render"
	"net/http"
)

// ErrorJSON makes json and respond with error
//
// If err is an instance of httpErr, then `httpStatusCode`, `details` and
// `errCode` are ignored and those fields are pass from httpErr instance.
func ErrorJSON(w http.ResponseWriter, r *http.Request, httpStatusCode int,
	err error, details string, errCode int) {
	errorMsg := ""

	switch e := err.(type) {
	case httpErr:
		if werr := e.Unwrap(); werr != nil {
			errorMsg = werr.Error()
		}
		httpStatusCode = e.statusCode
		errCode = e.errorCode
		details = e.msg
	case nil:
	default:
		errorMsg = err.Error()
	}

	render.Status(r, httpStatusCode)
	render.JSON(w, r, map[string]interface{}{
		"code":    errCode,
		"error":   errorMsg,
		"details": details})
}

// New5XXHttpErr returns http error that would be returned to client with
// 500 http status code and msg string ad details.
func New5XXHttpErr(err error, msg string) error {
	return httpErr{
		statusCode: http.StatusInternalServerError,
		msg:        msg,
		err:        err,
		errorCode:  0,
	}
}

type httpErr struct {
	statusCode int
	msg        string
	err        error
	errorCode  int
}

func (err httpErr) Error() string {
	return err.msg
}

func (err httpErr) Unwrap() error {
	return err.err
}
