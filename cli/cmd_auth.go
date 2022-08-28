package cli

import (
	"fmt"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"github.com/txsvc/stdlib/v2"

	"github.com/txsvc/apikit"
	"github.com/txsvc/apikit/api"
	"github.com/txsvc/apikit/config"
	"github.com/txsvc/apikit/helpers"
	"github.com/txsvc/apikit/internal/auth"
	"github.com/txsvc/apikit/internal/settings"
	"github.com/txsvc/apikit/logger"
)

func WithAuthCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "auth",
			Aliases: []string{"a"},
			Usage:   "options to register and login",
			Subcommands: []*cli.Command{
				{
					Name:        "init",
					Usage:       "register with the API service",
					UsageText:   "init email [passphrase]", // FIXME: better description
					Description: "longform description",    // FIXME: better description
					Action:      InitCommand,
				},
				{
					Name:        "login",
					Usage:       "authenticate with the API service",
					UsageText:   "login token",          // FIXME: better description
					Description: "longform description", // FIXME: better description
					Action:      LoginCommand,
				},
				{
					Name:        "logout",
					Usage:       "logout from the API service",
					UsageText:   "logout",               // FIXME: better description
					Description: "longform description", // FIXME: better description
					Action:      LogoutCommand,
				},
			},
		},
	}
}

func InitCommand(c *cli.Context) error {
	if c.NArg() < 1 || c.NArg() > 2 {
		return apikit.ErrInvalidNumArguments
	}

	userid := ""
	phrase := ""

	if c.NArg() == 1 {
		userid = c.Args().First() // path only
	} else if c.NArg() == 2 {
		userid = c.Args().First() // path and pass phrase
		phrase = c.Args().Get(1)
	}

	// create or validate the words
	mnemonic, err := helpers.CreateMnemonic(phrase)
	if err != nil {
		return err
	}

	// load settings
	cfg := config.GetSettings()

	// get a client instance
	cl, err := api.NewClient(cfg, logger.New())
	if err != nil {
		return err // FIXME: better err or just pass on what comes?
	}

	// if a passphrase was provided and the fingerprint(realm,userid,phrase) matches the API key,
	// then the user is re-initializing an existing account which is allowed. The client just
	// sends a logout request first before initiating the normal auth sequence.

	_apiKey := stdlib.Fingerprint(fmt.Sprintf("%s%s%s", config.AppInfo().Name(), userid, mnemonic))

	switch cfg.Status {
	case -1:
		// set to INVALID
		return config.ErrInvalidConfiguration
	case 1:
		if _apiKey == cfg.APIKey {
			// correct pass phrase was provided, reset the authentication
			if err := cl.LogoutCommand(); err != nil {
				return err // FIXME: better err or just pass on what comes?
			}
		} else {
			// already authenticated, abort
			return auth.ErrAlreadyAuthorized
		}
	}

	// 0, -2: don't care, can be overwritten as the client is not authorized yet

	cfg.Credentials = &settings.Credentials{
		ProjectID: config.AppInfo().Name(),
		UserID:    userid,
		Token:     api.CreateSimpleToken(),
		Expires:   0, // FIXME: should this expire after some time?
	}
	cfg.Status = settings.StateInit
	cfg.APIKey = _apiKey
	cfg.Scopes = make([]string, 0)
	cfg.DefaultScopes = make([]string, 0)
	cfg.Options = make(map[string]string)

	// now start the auth init process with the API

	err = cl.InitCommand(cfg)
	if err != nil {
		return err // FIXME: better err or just pass on what comes?
	}

	// finally save the file
	pathToFile := filepath.Join(config.ResolveConfigLocation(), config.DefaultConfigFileName)
	if err := helpers.WriteSettingsToFile(cfg, pathToFile); err != nil {
		return config.ErrInitializingConfiguration
	}

	if phrase == "" {
		fmt.Printf("userid: %s\n", cfg.Credentials.UserID)
		fmt.Printf("passphrase: \"%s\"\n\n", mnemonic)
		fmt.Println("Make a copy of the passphrase and keep it secure !")
	}

	return nil
}

func LoginCommand(c *cli.Context) error {
	if c.NArg() < 1 || c.NArg() > 1 {
		return apikit.ErrInvalidNumArguments
	}

	token := c.Args().First()

	// load settings
	cfg := config.GetSettings()
	if !cfg.Credentials.IsValid() {
		return config.ErrInvalidConfiguration
	}

	// now start the auth login process with the API
	cl, err := api.NewClient(cfg, logger.New())
	if err != nil {
		return err // FIXME: better err or just pass on what comes?
	}
	status, err := cl.LoginCommand(token)
	if err != nil {
		return err // FIXME: better err or just pass on what comes?
	}

	// update the local config
	cfg.Credentials.Token = status.Message
	cfg.Status = settings.StateAuthorized // LOGGED_IN
	if !cfg.Credentials.IsValid() {
		return config.ErrInvalidConfiguration
	}

	pathToFile := filepath.Join(config.ResolveConfigLocation(), config.DefaultConfigFileName)
	if err := helpers.WriteSettingsToFile(cfg, pathToFile); err != nil {
		return config.ErrInitializingConfiguration
	}

	fmt.Println("auth login done") // FIXME: better message !

	return nil
}

func LogoutCommand(c *cli.Context) error {
	if c.NArg() > 0 {
		return apikit.ErrInvalidNumArguments
	}

	// load settings
	cfg := config.GetSettings()
	if !cfg.Credentials.IsValid() {
		return config.ErrInvalidConfiguration
	}

	// now start the auth logout process with the API
	cl, err := api.NewClient(cfg, logger.New())
	if err != nil {
		return err // FIXME: better err or just pass on what comes?
	}
	err = cl.LogoutCommand()
	if err != nil {
		return err // FIXME: better err or just pass on what comes?
	}

	// update the local config
	cfg.Credentials.Expires = stdlib.Now() - 1
	cfg.Status = settings.StateUndefined // LOGGED_OUT

	pathToFile := filepath.Join(config.ResolveConfigLocation(), config.DefaultConfigFileName)
	if err := helpers.WriteSettingsToFile(cfg, pathToFile); err != nil {
		return config.ErrInitializingConfiguration
	}

	fmt.Println("auth logout done") // FIXME: better message !

	return nil
}
