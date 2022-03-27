package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDnsDetect(t *testing.T) {

	expected := "cloudflare"
	var expectedErr error = nil

	result, err := dnsDetect()

	assert.Equalf(t, expected, result, "Expected result: %s, got %s", expected, result)
	assert.Equalf(t, expectedErr, err, "Expected error: %s, got %s", expectedErr, err)
}
