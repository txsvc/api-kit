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
	"github.com/txsvc/apikit/internal"
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

	// if a passphrase was provided and the fingerprint(realm,userid,phrase) matches the API key,
	// then the user is re-initializing an existing account which is allowed. The client just
	// sends a logout request first before initiating the normal auth sequence.

	_apiKey := stdlib.Fingerprint(fmt.Sprintf("%s%s%s", config.Name(), userid, mnemonic))

	if cfg.Status != -2 {
		if cfg.APIKey == _apiKey {
			// FIXME: send logout
			fmt.Println("logout")
		} else {
			if cfg.Credentials.Token != "" {
				return auth.ErrAlreadyInitialized // FIXME: can we do better ?
			}
		}
	}

	cfg.Credentials = &settings.Credentials{
		ProjectID: config.Name(),
		UserID:    userid,
		Token:     internal.CreateSimpleToken(),
		Expires:   0, // FIXME: should this expire after some time?
	}
	cfg.Status = -2
	cfg.APIKey = _apiKey

	// now start the auth init process with the API
	cl, err := api.NewClient(cfg, logger.New())
	if err != nil {
		return err // FIXME: better err or just pass on what comes?
	}
	err = cl.InitCommand(cfg)
	if err != nil {
		return err // FIXME: better err or just pass on what comes?
	}

	// finally save the file
	pathToFile := filepath.Join(config.ResolveConfigLocation(), config.DefaultConfigFileName)
	if err := cfg.WriteToFile(pathToFile); err != nil {
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
	cfg.Status = 1 // LOGGED_IN
	if !cfg.Credentials.IsValid() {
		return config.ErrInvalidConfiguration
	}

	pathToFile := filepath.Join(config.ResolveConfigLocation(), config.DefaultConfigFileName)
	if err := cfg.WriteToFile(pathToFile); err != nil {
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
	cfg.Status = -1 // LOGGED_OUT

	pathToFile := filepath.Join(config.ResolveConfigLocation(), config.DefaultConfigFileName)
	if err := cfg.WriteToFile(pathToFile); err != nil {
		return config.ErrInitializingConfiguration
	}

	fmt.Println("auth logout done") // FIXME: better message !

	return nil
}
