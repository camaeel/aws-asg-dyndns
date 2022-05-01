package cloudflare_dns

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/camaeell/aws-asg-dyndns/awsClient"
	"github.com/cloudflare/cloudflare-go"
)

type cloudflareProvider struct {
	token string
	zone  string
	api   CLOUDFLAREAPI
}

func NewCloudflareProvider(ctx context.Context, ssmClient awsClient.SSMAPI, domain string) (*cloudflareProvider, error) {
	ret := cloudflareProvider{}
	ret.detectZone(domain)
	token, err := awsClient.GetSSMParameterValue(ctx, ssmClient, ssmParameterTokenPath(ret.zone))
	if err != nil {
		return nil, err
	}
	ret.token = token

	api, err := cloudflare.NewWithAPIToken(ret.token)
	if err != nil {
		return nil, err
	}
	ret.api = api

	return &ret, nil
}

// Temporary solution
// TODO Implement proper solution.
func (c *cloudflareProvider) detectZone(domain string) {
	splitted := strings.Split(domain, ".")
	zone := strings.Join(splitted[len(splitted)-2:], ".")
	c.zone = zone
}

func ssmParameterTokenPath(zone string) string {
	return fmt.Sprintf("/dyn-dns/%s/cloudflare/token", zone)
}

func (c cloudflareProvider) DnsEntryAddIp(ctx context.Context, domain string, ip *string) error {

	zoneId, err := c.api.ZoneIDByName(c.zone)
	if err != nil {
		return err
	}

	dnsRecordQuery := cloudflare.DNSRecord{Name: domain, Type: "A", Content: *ip}

	dnsRecords, err := c.api.DNSRecords(ctx, zoneId, dnsRecordQuery)
	if err != nil {
		return err
	}

	if len(dnsRecords) > 1 {
		return fmt.Errorf("Found too many DNS records: %d for %s domain and ip = %s", len(dnsRecords), domain, *ip)
	} else if len(dnsRecords) == 1 {
		log.Printf("Warning. DNS records already exists for %s domain and ip = %s", domain, *ip)
	} else {
		dnsRecord := cloudflare.DNSRecord{Name: domain, Type: "A", Content: *ip, TTL: 60}
		_, err := c.api.CreateDNSRecord(ctx, zoneId, dnsRecord)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c cloudflareProvider) DnsEntryRemoveIp(ctx context.Context, domain string, ip *string) error {
	zoneId, err := c.api.ZoneIDByName(c.zone)
	if err != nil {
		return err
	}

	dnsRecordQuery := cloudflare.DNSRecord{Name: domain, Type: "A", Content: *ip}

	dnsRecords, err := c.api.DNSRecords(ctx, zoneId, dnsRecordQuery)
	if err != nil {
		return err
	}

	if len(dnsRecords) > 1 {
		return fmt.Errorf("Found too many (%d) DNS records for %s domain and %s", len(dnsRecords), domain, *ip)
	} else if len(dnsRecords) == 1 {
		err := c.api.DeleteDNSRecord(ctx, zoneId, dnsRecords[0].ID)
		if err != nil {
			return err
		}
	} else {
		log.Printf("Warning. DNS records already doesn't exist for %s domain and ip = %s", domain, *ip)
	}

	return nil
}
