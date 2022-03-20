package dns

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type DnsProvider interface {
	dnsEntryAddIp(ctx context.Context, domain string, ip *string) error
	dnsEntryRemoveIp(ctx context.Context, domain string, ip *string) error
}

func createDnsProvider(ctx context.Context, ssmClient *ssm.Client, dnsProviderName string, domain string) (DnsProvider, error) {
	if dnsProviderName == "cloudflare" {
		ret, err := newCloudflareProvider(ctx, ssmClient, domain)
		return ret, err
	}

	return nil, errors.New("Uknown DNS provider: " + dnsProviderName)
}
