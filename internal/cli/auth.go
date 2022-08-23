package cli

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"github.com/txsvc/apikit"
	"github.com/txsvc/apikit/config"
	"github.com/txsvc/apikit/internal"
)

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

	// FIXME: add a flag to overwrite the existing config

	cfg := config.GetSettings()
	if cfg.Credentials != nil {
		return apikit.ErrAlreadyInitialized // FIXME: can we do better ?
	}

	mnemonic, err := internal.CreateMnemonic(phrase)
	if err != nil {
		return err
	}

	// create credentials and keys
	_cfg, err := internal.InitSettings(config.Name(), userid)
	if err != nil {
		log.Fatal(err)
	}

	// copy the credentials and api keys
	cfg.Credentials = _cfg.Credentials
	cfg.APIKey = _cfg.APIKey

	pathToFile := filepath.Join(config.ResolveConfigLocation(), config.DefaultConfigFileName)
	if err := cfg.WriteToFile(pathToFile); err != nil {
		return apikit.ErrInitializingConfiguration
	}

	if phrase == "" {
		fmt.Printf("user-id: %s\n", cfg.Credentials.UserID)
		fmt.Printf("passphrase: \"%s\"\n\n", mnemonic)
		fmt.Println("Make a copy of the pass phrase and keep it secure !")
	}

	return nil
}
