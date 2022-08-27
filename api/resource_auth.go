package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/txsvc/apikit/config"
	"github.com/txsvc/apikit/helpers"
	"github.com/txsvc/apikit/internal"
	"github.com/txsvc/apikit/internal/auth"
	"github.com/txsvc/apikit/internal/settings"
	"github.com/txsvc/stdlib/v2"
)

const (
	// auth routes
	InitRoute   = "/auth"
	LoginRoute  = "/auth/:sig/:token"
	LogoutRoute = "/auth/:sig"
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

func (c *Client) InitCommand(cfg *settings.Settings) error {
	_, err := c.POST(fmt.Sprintf("%s%s", NamespacePrefix, InitRoute), cfg, nil)
	return err
}

func InitEndpoint(c echo.Context) error {
	// get the payload
	var cfg *settings.Settings = new(settings.Settings)
	if err := c.Bind(cfg); err != nil {
		return StandardResponse(c, http.StatusBadRequest, nil)
	}

	// pre-validate the request
	if cfg.Credentials == nil || cfg.APIKey == "" {
		return StandardResponse(c, http.StatusBadRequest, nil)
	}
	if cfg.Credentials.ProjectID == "" || cfg.Credentials.UserID == "" {
		return StandardResponse(c, http.StatusBadRequest, nil)
	}

	// prepare the settings for registration
	cfg.Credentials.Token = internal.CreateSimpleToken()    // ignore anything that was provided
	cfg.Credentials.Expires = stdlib.IncT(stdlib.Now(), 15) // FIXME: config, valid for 15min
	cfg.Status = -2                                         // signals init

	if err := auth.RegisterAuthorization(cfg); err != nil {
		return StandardResponse(c, http.StatusBadRequest, nil) // FIXME: or 409/Conflict ?
	}

	// all good so far, send the confirmation
	err := helpers.MailgunSimpleEmail("ops@txs.vc", cfg.Credentials.UserID, "auth", fmt.Sprintf("the token: %s\n", cfg.Credentials.Token))
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
	_cfg, err := auth.LookupByToken(token)
	if _cfg == nil && err != nil {
		return ErrorResponse(c, http.StatusBadRequest, ErrInternalError, "token")
	}
	if _cfg == nil && err == nil {
		return ErrorResponse(c, http.StatusBadRequest, config.ErrInitializingConfiguration, "not found") // simply not there ...
	}

	// compare provided signature with the expected signature
	if sig != signature(_cfg.APIKey, _cfg.Credentials.Token) {
		return ErrorResponse(c, http.StatusBadRequest, config.ErrInitializingConfiguration, "invalid sig")
	}

	// everything checks out, create/register the real credentials now ...
	cfg := _cfg.Clone()         // clone, otherwise stupid things happen with pointers !
	cfg.Credentials.Expires = 0 // FIXME: really never ?
	cfg.Credentials.Token = internal.CreateSimpleToken()
	cfg.Status = 1 // FIXME: LOGGED_IN as const

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
	cfg.Status = -1 // just set to invalid and expired
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
