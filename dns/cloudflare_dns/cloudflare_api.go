package cloudflare_dns

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
)

type CLOUDFLAREAPI interface {
	DeleteDNSRecord(ctx context.Context, zoneID string, recordID string) error
	DNSRecords(ctx context.Context, zoneID string, rr cloudflare.DNSRecord) ([]cloudflare.DNSRecord, error)
	ZoneIDByName(zoneName string) (string, error)
	CreateDNSRecord(ctx context.Context, zoneID string, rr cloudflare.DNSRecord) (*cloudflare.DNSRecordResponse, error)
}
