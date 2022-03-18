package dns

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// Adds IP address to DNS entry
func DnsEntryAddIp(ctx context.Context, ssmClient *ssm.Client, domain string, ip *string) error {
	dnsProviderName, err := dnsDetect()
	if err != nil {
		return err
	}

	dnsProvider, err := createDnsProvider(ctx, ssmClient, dnsProviderName, domain)
	if err != nil {
		return err
	}

	return dnsProvider.dnsEntryAddIp(domain, ip)
}

// Removes IP address from DNS entry
func DnsEntryRemoveIp(ctx context.Context, ssmClient *ssm.Client, domain string, ip *string) error {
	dnsProviderName, err := dnsDetect()
	if err != nil {
		return err
	}

	dnsProvider, err := createDnsProvider(ctx, ssmClient, dnsProviderName, domain)
	if err != nil {
		return err
	}

	return dnsProvider.dnsEntryRemoveIp(domain, ip)
}

func dnsDetect() (string, error) {
	// TODO implement logic
	// Detect from /dyn-dns/ZONE/provider
	return "cloudflare", nil
}
