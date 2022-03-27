package dns

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/camaeell/aws-asg-dyndns/dns/cloudflare_dns"
)

type DnsProvider interface {
	DnsEntryAddIp(ctx context.Context, domain string, ip *string) error
	DnsEntryRemoveIp(ctx context.Context, domain string, ip *string) error
}

func createDnsProvider(ctx context.Context, ssmClient *ssm.Client, dnsProviderName string, domain string) (DnsProvider, error) {
	if dnsProviderName == "cloudflare" {
		ret, err := cloudflare_dns.NewCloudflareProvider(ctx, ssmClient, domain)
		return ret, err
	}

	return nil, errors.New("Uknown DNS provider: " + dnsProviderName)
}

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

	return dnsProvider.DnsEntryAddIp(ctx, domain, ip)
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

	return dnsProvider.DnsEntryRemoveIp(ctx, domain, ip)
}
