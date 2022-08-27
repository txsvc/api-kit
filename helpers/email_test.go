package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/txsvc/stdlib/v2"
)

func TestMailgunSimpleEmail(t *testing.T) {
	domain := stdlib.GetString(MailgunEmailDomainENV, "")
	apiKey := stdlib.GetString(MailgunApiKeyENV, "")
	if domain == "" || apiKey == "" {
		return // skip the test
	}

	err := MailgunSimpleEmail("test@podops.dev", "ops@txs.vc", "TestMailgunSimpleEmail", "Testing TestMailgunSimpleEmail")
	assert.NoError(t, err)
}
