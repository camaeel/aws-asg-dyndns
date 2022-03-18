package dns

import (
	"testing"

	"gotest.tools/assert"
)

func TestCloudflareDomainToZoneResolve(t *testing.T) {
	domain := "test.test1.example.com"
	expected := "example.com"

	var cfp cloudflareProvider
	cfp.detectZone(domain)

	assert.Equal(t, expected, cfp.zone, "Zone should be equal.")
}
