package internal

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"

	"github.com/txsvc/apikit/internal/settings"
)

const (
	// default API scopes
	ScopeApiRead   = "api:read"
	ScopeApiWrite  = "api:write"
	ScopeApiEdit   = "api:edit"
	ScopeApiCreate = "api:create"
	ScopeApiDelete = "api:delete"
	ScopeApiAdmin  = "api:admin"
)

var (
	// ErrNotAuthorized indicates that the API caller is not authorized
	ErrNotAuthorized     = errors.New("not authorized")
	ErrAlreadyAuthorized = errors.New("already authorized")

	// ErrNoToken indicates that no bearer token was provided
	ErrNoToken = errors.New("no token provided")
	// ErrNoScope indicates that no scope was provided
	ErrNoScope = errors.New("no scope provided")

	// different types of lookup tables
	tokenToAuth map[string]*settings.Settings
	idToAuth    map[string]*settings.Settings
	mu          sync.Mutex // used to protect the above maps
)

func init() {
	// just empty maps to avoid any NPEs
	FlushAuthorizations("")
}

// CheckAuthorization relies on the presence of a bearer token and validates the
// matching authorization against a list of requested scopes. If everything checks out,
// the function returns the authorization or an error otherwise.
func CheckAuthorization(ctx context.Context, c echo.Context, scope string) (*settings.Settings, error) {
	token, err := GetBearerToken(c.Request())
	if err != nil {
		return nil, err
	}

	auth, err := FindAuthorizationByToken(ctx, token)
	if err != nil || auth == nil {
		return nil, ErrNotAuthorized
	}

	if hasScope(auth.GetScopes(), ScopeApiAdmin) {
		return auth, nil
	}
	if !hasScope(auth.GetScopes(), scope) {
		return nil, ErrNotAuthorized
	}

	return auth, nil
}

func GetBearerToken(r *http.Request) (string, error) {

	// FIXME: optimize this !!

	auth := r.Header.Get("Authorization")
	if len(auth) == 0 {
		return "", ErrNoToken
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 {
		return "", ErrNoToken
	}
	if parts[0] == "Bearer" {
		return parts[1], nil
	}

	return "", ErrNoToken
}

func FlushAuthorizations(root string) {
	mu.Lock()
	defer mu.Unlock()

	tokenToAuth = make(map[string]*settings.Settings)
	idToAuth = make(map[string]*settings.Settings)

	if root != "" {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				cfg, err := ReadSettingsFromFile(path)
				if err != nil {
					return err // FIXME: this is never checked on exit !
				}
				RegisterAuthorization(cfg)
			}
			return nil
		})
	}
}

func RegisterAuthorization(cfg *settings.Settings) {
	tokenToAuth[cfg.Credentials.Token] = cfg
	idToAuth[key(cfg.Credentials.ProjectID, cfg.Credentials.UserID)] = cfg
}

func LookupAuthorization(ctx context.Context, realm, userid string) (*settings.Settings, error) {
	if a, ok := idToAuth[key(realm, userid)]; ok {
		return a, nil
	}
	return nil, nil // FIXME: return an error ?
}

func FindAuthorizationByToken(ctx context.Context, token string) (*settings.Settings, error) {
	if token == "" {
		return nil, ErrNoToken
	}
	if a, ok := tokenToAuth[token]; ok {
		return a, nil
	}
	return nil, nil // FIXME: return an error ?
}
