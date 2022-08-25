package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/txsvc/stdlib/v2"

	"github.com/txsvc/apikit/helpers"
	"github.com/txsvc/apikit/internal/settings"
)

func CreateSimpleToken() string {
	token, _ := stdlib.UUID()
	return token
}

func ReadSettingsFromFile(path string) (*settings.Settings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ds := settings.Settings{}
	if err := json.Unmarshal([]byte(data), &ds); err != nil {
		return nil, err
	}
	return &ds, nil
}

func ReadCredentialsFromFile(path string) (*settings.Credentials, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cred := settings.Credentials{}
	if err := json.Unmarshal([]byte(data), &cred); err != nil {
		return nil, err
	}
	return &cred, nil
}

func InitSettings(realm, userid string) (*settings.Settings, error) {
	cred := InitCredentials(realm, userid)
	cfg := settings.Settings{
		Credentials: &cred,
	}

	// create a mnemic to derieve the 'password' from
	mnemonic, err := helpers.CreateMnemonic("")
	if err != nil {
		return nil, err // abort here
	}
	cfg.APIKey = stdlib.Fingerprint(fmt.Sprintf("%s%s%s", realm, userid, mnemonic))
	cfg.SetOption("PASSPHRASE", mnemonic) // FIXME: make sure this gets NEVER saved to disc!

	return &cfg, nil
}

func InitCredentials(realm, userid string) settings.Credentials {
	return settings.Credentials{
		ProjectID: realm,
		UserID:    userid,
		Token:     CreateSimpleToken(),
		Expires:   0,
	}
}

func key(part1, part2 string) string {
	return strings.ToLower(part1 + "." + part2)
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
