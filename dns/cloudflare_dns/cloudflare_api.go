package cloudflare_dns

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
)

//go:generate mockgen -destination=./mocks/aws_apis_mocks.go -package=mocks github.com/camaeell/aws-asg-dyndns/dns/cloudflare_dns CLOUDFLAREAPI

type CLOUDFLAREAPI interface {
	DeleteDNSRecord(ctx context.Context, zoneID string, recordID string) error
	DNSRecords(ctx context.Context, zoneID string, rr cloudflare.DNSRecord) ([]cloudflare.DNSRecord, error)
	ZoneIDByName(zoneName string) (string, error)
	CreateDNSRecord(ctx context.Context, zoneID string, rr cloudflare.DNSRecord) (*cloudflare.DNSRecordResponse, error)
}
