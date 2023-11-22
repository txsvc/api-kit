package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/txsvc/cloudlib/settings"
)

const (
	// anonymous
	ScopeAnonymous = "api:anonymous" // this basically means that the Client is unknown
	// default API scopes
	ScopeApiRead   = "api:read" // that's the very minimum for a proper client
	ScopeApiWrite  = "api:write"
	ScopeApiEdit   = "api:edit"
	ScopeApiCreate = "api:create"
	ScopeApiDelete = "api:delete"
	ScopeApiAdmin  = "api:admin"
	// block access
	ScopeApiNoAccess = "api:noaccess"
)

var (
	// ErrNotAuthorized indicates that the API caller is not authorized
	ErrNotAuthorized     = errors.New("not authorized")
	ErrAlreadyAuthorized = errors.New("already authorized")

	// ErrAlreadyInitialized indicates that client is already registered
	ErrAlreadyInitialized = errors.New("already initialized")

	// ErrNoToken indicates that no bearer token was provided
	ErrNoToken = errors.New("no token provided")
	// ErrTokenExpired indicates that the token is no longer valid
	ErrTokenExpired = errors.New("token expired")

	// ErrNoScope indicates that no scope was provided
	ErrNoScope = errors.New("no scope provided")
)

func init() {
	// just empty maps to avoid any NPEs
	FlushAuthorizations("")
}

// CheckAuthorization relies on the presence of a bearer token and validates the
// matching authorization against a list of requested scopes. If everything checks out,
// the function returns the authorization or an error otherwise.
func CheckAuthorization(ctx context.Context, c echo.Context, scope string) (*settings.DialSettings, error) {
	token, err := GetBearerToken(c.Request())
	if err != nil {
		return nil, err
	}

	auth, err := cache.LookupByToken(token)
	if err != nil || auth == nil || !auth.Credentials.IsValid() {
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

// FIXME: this is a VERY simple implementation
func hasScope(target []string, scope string) bool {

	scopes := strings.Split(scope, ",")
	mustMatch := len(scopes)

	for _, s := range scopes {
		for _, ss := range target {
			if s == ss {
				mustMatch--
				break
			}
		}
	}

	return mustMatch == 0
}
