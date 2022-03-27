package cloudflare_dns

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
)

type CloudflareFakeClient struct {
	DeleteDnsRecordErr      error
	DNSRecordsResponse      []cloudflare.DNSRecord
	DNSRecordsErr           error
	ZoneIdByNameResponse    string
	ZoneIdByNameErr         error
	CreateDNSRecordResponse *cloudflare.DNSRecordResponse
	CreateDNSRecordErr      error
}

func (c *CloudflareFakeClient) DeleteDNSRecord(ctx context.Context, zoneID string, recordID string) error {
	return c.DeleteDnsRecordErr
}

func (c *CloudflareFakeClient) DNSRecords(ctx context.Context, zoneID string, rr cloudflare.DNSRecord) ([]cloudflare.DNSRecord, error) {
	return c.DNSRecordsResponse, c.DNSRecordsErr
}

func (c *CloudflareFakeClient) ZoneIDByName(zoneName string) (string, error) {
	return c.ZoneIdByNameResponse, c.ZoneIdByNameErr

}

func (c *CloudflareFakeClient) CreateDNSRecord(ctx context.Context, zoneID string, rr cloudflare.DNSRecord) (*cloudflare.DNSRecordResponse, error) {
	return c.CreateDNSRecordResponse, c.CreateDNSRecordErr
}
