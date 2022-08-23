package internal

import (
	"strings"

	"github.com/Bytom/bytom/wallet/mnemonic"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"

	"github.com/txsvc/apikit"
)

const (
	MinWordsInPassPhrase = 11
)

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
		return "", apikit.ErrInvalidPassPhrase
	}

	return mnemonicPhrase, nil
}
