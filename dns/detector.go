package dns

import "github.com/aws/aws-sdk-go-v2/service/ssm"

type DnsProvider interface {
	dnsEntryAddIp(domain string, ip *string) error
	dnsEntryRemoveIp(domain string, ip *string) error
}

// Adds IP address to DNS entry
func DnsEntryAddIp(ssmClient *ssm.Client, domain string, ip *string) error {

	dnsProvider, err := dnsDetect()
	if err != nil {
		return err
	}
	return dnsProvider.dnsEntryAddIp(domain, ip)
}

// Removes IP address from DNS entry
func DnsEntryRemoveIp(ssmClient *ssm.Client, domain string, ip *string) error {
	dnsProvider, err := dnsDetect()
	if err != nil {
		return err
	}
	return dnsProvider.dnsEntryRemoveIp(domain, ip)
}

func dnsDetect() (DnsProvider, error) {
	cf := cloudflare{}
	// TODO implement logic
	return cf, nil
}
