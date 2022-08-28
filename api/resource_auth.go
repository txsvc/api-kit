package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/txsvc/stdlib/v2"

	"github.com/txsvc/apikit/config"
	"github.com/txsvc/apikit/helpers"
	"github.com/txsvc/apikit/internal/auth"
	"github.com/txsvc/apikit/internal/settings"
)

const (
	// auth routes
	InitRoute   = "/auth"
	LoginRoute  = "/auth/:sig/:token"
	LogoutRoute = "/auth/:sig"

	LoginExpiresAfter = 15
)

func WithAuthEndpoints(e *echo.Echo) *echo.Echo {
	// grouped under /a/v1
	apiGroup := e.Group(NamespacePrefix)

	// add the routes
	apiGroup.POST(InitRoute, InitEndpoint)
	apiGroup.GET(LoginRoute, LoginEndpoint)
	apiGroup.DELETE(LogoutRoute, LogoutEndpoint)

	// done
	return e
}

func (c *Client) InitCommand(cfg *settings.DialSettings) error {
	_, err := c.POST(fmt.Sprintf("%s%s", NamespacePrefix, InitRoute), cfg, nil)
	return err
}

func InitEndpoint(c echo.Context) error {
	// get the payload
	var cfg_ *settings.DialSettings = new(settings.DialSettings)
	if err := c.Bind(cfg_); err != nil {
		return StandardResponse(c, http.StatusBadRequest, nil)
	}

	// pre-validate the request
	if cfg_.Credentials == nil || cfg_.APIKey == "" {
		return StandardResponse(c, http.StatusBadRequest, nil)
	}
	if cfg_.Credentials.ProjectID == "" || cfg_.Credentials.UserID == "" {
		return StandardResponse(c, http.StatusBadRequest, nil)
	}

	// create a brand new instance so that the client can't sneak anything in we don't want
	cfg := settings.DialSettings{
		Credentials:   cfg_.Credentials.Clone(),
		DefaultScopes: config.GetConfig().GetScopes(),
	}

	// prepare the settings for registration
	cfg.Credentials.Token = CreateSimpleToken() // ignore anything that was provided
	cfg.Credentials.Expires = stdlib.IncT(stdlib.Now(), LoginExpiresAfter)
	cfg.APIKey = cfg_.APIKey
	cfg.Status = settings.StateInit // signals init

	if err := auth.RegisterAuthorization(&cfg); err != nil {
		return StandardResponse(c, http.StatusBadRequest, nil) // FIXME: or 409/Conflict ?
	}

	// all good so far, send the confirmation
	err := helpers.MailgunSimpleEmail("ops@txs.vc", cfg.Credentials.UserID, fmt.Sprintf("your api access credentials (%d)", stdlib.Now()), fmt.Sprintf("the token: %s\n", cfg.Credentials.Token))
	if err != nil {
		return StandardResponse(c, http.StatusBadRequest, nil)
	}
	// FIXME: the email sending has to be better !

	return StandardResponse(c, http.StatusCreated, nil)
}

func (c *Client) LoginCommand(token string) (*StatusObject, error) {
	var so StatusObject

	status, err := c.GET(fmt.Sprintf("%s%s/%s/%s", NamespacePrefix, InitRoute, signature(c.cfg.APIKey, token), token), &so)
	if status != http.StatusOK || err != nil {
		return nil, err
	}
	return &so, nil
}

func LoginEndpoint(c echo.Context) error {
	sig := c.Param("sig")
	if sig == "" {
		return ErrorResponse(c, http.StatusBadRequest, ErrInvalidRoute, "sig")
	}
	token := c.Param("token")
	if token == "" {
		return ErrorResponse(c, http.StatusBadRequest, ErrInvalidRoute, "token")
	}

	// verify the request
	cfg_, err := auth.LookupByToken(token)
	if cfg_ == nil && err != nil {
		return ErrorResponse(c, http.StatusBadRequest, ErrInternalError, "token")
	}
	if cfg_ == nil && err == nil {
		return ErrorResponse(c, http.StatusBadRequest, config.ErrInitializingConfiguration, "not found") // simply not there ...
	}

	// compare provided signature with the expected signature
	if sig != signature(cfg_.APIKey, cfg_.Credentials.Token) {
		return ErrorResponse(c, http.StatusBadRequest, config.ErrInitializingConfiguration, "invalid sig")
	}

	// check if the token is still valid
	if cfg_.Credentials.Expires < stdlib.Now() {
		return ErrorResponse(c, http.StatusBadRequest, auth.ErrTokenExpired, "expired")
	}

	// everything checks out, create/register the real credentials now ...
	cfg := cfg_.Clone()         // clone, otherwise stupid things happen with pointers !
	cfg.Credentials.Expires = 0 // FIXME: really never ?
	cfg.Credentials.Token = CreateSimpleToken()
	cfg.Status = settings.StateAuthorized

	// FIXME: what about scopes ?

	if err := auth.RegisterAuthorization(&cfg); err != nil {
		fmt.Println(err)
		return ErrorResponse(c, http.StatusBadRequest, config.ErrInitializingConfiguration, "can't register")
	}

	// just send the token back
	resp := StatusObject{
		Status:  http.StatusOK,
		Message: cfg.Credentials.Token,
	}

	return StandardResponse(c, http.StatusOK, resp)
}

func (c *Client) LogoutCommand() error {
	_, err := c.DELETE(fmt.Sprintf("%s%s/%s", NamespacePrefix, InitRoute, signature(c.cfg.APIKey, c.cfg.Credentials.Token)), nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func LogoutEndpoint(c echo.Context) error {
	sig := c.Param("sig")
	if sig == "" {
		return ErrorResponse(c, http.StatusBadRequest, ErrInvalidRoute, "sig")
	}
	token, err := auth.GetBearerToken(c.Request())
	if err != nil {
		return ErrorResponse(c, http.StatusUnauthorized, err, "")
	}

	// verify the request
	cfg, err := auth.LookupByToken(token)
	if cfg == nil && err != nil {
		return ErrorResponse(c, http.StatusBadRequest, ErrInternalError, "token")
	}

	// compare provided signature with the expected signature
	if sig != signature(cfg.APIKey, cfg.Credentials.Token) {
		return ErrorResponse(c, http.StatusBadRequest, config.ErrInitializingConfiguration, "invalid sig")
	}

	// update the cache and store
	cfg.Status = settings.StateUndefined // just set to invalid and expired
	cfg.Credentials.Expires = stdlib.Now() - 1
	if err := auth.UpdateStore(cfg); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err, "update store")
	}

	return StandardResponse(c, http.StatusOK, nil)
}

// signature returns a MD5(apiKey+token) as this is only known locally ...
func signature(apiKey, token string) string {
	return stdlib.Fingerprint(fmt.Sprintf("%s%s", apiKey, token))
}

func CreateSimpleToken() string {
	token, _ := stdlib.UUID()
	return token
}
