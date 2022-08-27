package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type (
	// StatusObject is used to report operation status and errors in an API request.
	// The struct can be used as a response object or be treated as an error object
	StatusObject struct {
		Status    int    `json:"status" binding:"required"`
		Message   string `json:"message" binding:"required"`
		RootError error  `json:"-"`
	}

	// RelevantHeaders represents the most important headers
	RelevantHeaders struct {
		Range           string `header:"Range"`
		UserAgent       string `header:"User-Agent"`
		Forwarded       string `header:"Forwarded"`
		XForwardedFor   string `header:"X-Forwarded-For"`
		XForwwardedHost string `header:"X-Forwarded-Host"`
		Referer         string `header:"Referer"`
	}
)

var (
	// ErrInvalidRoute indicates that the route and/or its parameters are not valid
	ErrInvalidRoute = errors.New("invalid route")
	// ErrNotImplemented indicates that a function is not yet implemented
	ErrNotImplemented = errors.New("not implemented")
	// ErrInternalError indicates everything else
	ErrInternalError = errors.New("internal error")
)

func (h *RelevantHeaders) Ranges() (int64, int64) {
	return ParseRange(h.Range)
}

// NewStatus initializes a new StatusObject
func NewStatus(s int, m string) StatusObject {
	return StatusObject{Status: s, Message: m}
}

// NewErrorStatus initializes a new StatusObject from an error
func NewErrorStatus(s int, e error, hint string) StatusObject {
	return StatusObject{Status: s, Message: fmt.Sprintf("%s (%s)", e.Error(), hint), RootError: e}
}

func (so *StatusObject) String() string {
	return fmt.Sprintf("%s: %d", so.Message, so.Status)
}

func (so *StatusObject) Error() string {
	return so.String()
}

// StandardResponse is the default way to respond to API requests
func StandardResponse(c echo.Context, status int, res interface{}) error {
	if res == nil {
		resp := StatusObject{
			Status:  status,
			Message: fmt.Sprintf("status: %d", status),
		}
		return c.JSON(status, &resp)
	} else {
		return c.JSON(status, res)
	}
}

// ErrorResponse reports the error and responds with an ErrorObject
func ErrorResponse(c echo.Context, status int, err error, hint string) error {
	var resp StatusObject

	// send the error to the Error Reporting
	// FIXME observer.ReportError(err)

	if err == nil {
		resp = NewStatus(http.StatusInternalServerError, fmt.Sprintf("%d", status))
	} else {
		resp = NewErrorStatus(status, err, hint)
	}
	return c.JSON(status, &resp)
}

// DefaultEndpoint just returns http.StatusOK
func DefaultEndpoint(c echo.Context) error {
	return StandardResponse(c, http.StatusOK, nil)
}

// HandleFileUpload receives files from stores them locally
func HandleFileUpload(ctx context.Context, req *http.Request, location, formName string) (string, error) {
	var path string

	// FIXME: treat location as a 'bucket' in preparation of switching to a generic storage API

	mr, err := req.MultipartReader()
	if err != nil {
		return "", err
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		if part.FormName() == formName {
			path = filepath.Join(location, part.FileName())

			os.MkdirAll(filepath.Dir(path), os.ModePerm) // make sure sub-folders exist
			out, err := os.Create(path)
			if err != nil {
				return "", err
			}
			defer out.Close()

			if _, err := io.Copy(out, part); err != nil {
				return "", err
			}
		}
	}

	return path, nil // FIXME: do we need the path ?
}
