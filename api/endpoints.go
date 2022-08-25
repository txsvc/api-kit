package api

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

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
