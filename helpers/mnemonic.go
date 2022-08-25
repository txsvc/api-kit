package helpers

import (
	"errors"
	"strings"

	"github.com/Bytom/bytom/wallet/mnemonic"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

const (
	MinWordsInPassPhrase = 11
)

// ErrInvalidPassPhrase indicates that the pass phrase is too short
var ErrInvalidPassPhrase = errors.New("invalid pass phrase")

func CreateMnemonic(phrase string) (string, error) {
	mnemonicPhrase := ""

	if phrase != "" {
		// make sure that each word is spaced with just one whitespace
		parts := strings.Split(phrase, " ")
		for _, s := range parts {
			if s != "" {
				mnemonicPhrase = mnemonicPhrase + " " + strings.Trim(s, " ")
			}
		}
		mnemonicPhrase = strings.Trim(mnemonicPhrase, " ")
	} else {
		seed, err := mnemonic.NewEntropy(128)
		if err != nil {
			return "", err
		}
		mnemonicPhrase, err = hdwallet.NewMnemonicFromEntropy(seed)
		if err != nil {
			return "", err
		}
	}

	if strings.Count(mnemonicPhrase, " ") < MinWordsInPassPhrase {
		return "", ErrInvalidPassPhrase
	}

	return mnemonicPhrase, nil
}
