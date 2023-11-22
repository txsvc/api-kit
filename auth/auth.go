package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/txsvc/cloudlib"
	"github.com/txsvc/cloudlib/settings"
)

const (
	TypeAuthProvider cloudlib.ProviderType = 40

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

type (
	AuthProvider interface {
		LookupByToken(token string) (*settings.DialSettings, error)
		UpdateStore(ds *settings.DialSettings) error
	}
)

var (
	// ErrInternalAuthError indicates that soemthing went wrong with the provider
	ErrInternalAuthError = errors.New("internal auth error")
	// ErrInvalidCredentials indicates that the provided credentials did not pass validation
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrTokenNotFound indicates that the token is not in the store
	ErrTokenNotFound = errors.New("token not found")

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

	authProvider *cloudlib.Provider
)

//
// The generic AuthProvider parts
//

func NewConfig(opts cloudlib.ProviderConfig) (*cloudlib.Provider, error) {
	if opts.Type != TypeAuthProvider {
		return nil, fmt.Errorf(cloudlib.MsgUnsupportedProviderType, opts.Type)
	}

	o, err := cloudlib.New(opts)
	if err != nil {
		return nil, err
	}
	authProvider = o

	return o, nil
}

func UpdateConfig(opts cloudlib.ProviderConfig) (*cloudlib.Provider, error) {
	if opts.Type != TypeAuthProvider {
		return nil, fmt.Errorf(cloudlib.MsgUnsupportedProviderType, opts.Type)
	}

	return authProvider, authProvider.RegisterProviders(true, opts)
}

func LookupByToken(token string) (*settings.DialSettings, error) {
	imp, found := authProvider.Find(TypeAuthProvider)
	if !found {
		return nil, ErrInternalAuthError
	}

	return imp.(AuthProvider).LookupByToken(token)
}

func UpdateStore(ds *settings.DialSettings) error {
	imp, found := authProvider.Find(TypeAuthProvider)
	if !found {
		return ErrInternalAuthError
	}

	return imp.(AuthProvider).UpdateStore(ds)
}

// Auth functionallity
//
// CheckAuthorization relies on the presence of a bearer token and validates the
// matching authorization against a list of requested scopes. If everything checks out,
// the function returns the authorization or an error otherwise.
func CheckAuthorization(ctx context.Context, c echo.Context, scope string) (*settings.DialSettings, error) {
	token, err := GetBearerToken(c.Request())
	if err != nil {
		return nil, err
	}

	auth, err := LookupByToken(token)
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
