package cloudflare_dns

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	awsClient_mocks "github.com/camaeell/aws-asg-dyndns/awsClient/mocks"
	cf_mocks "github.com/camaeell/aws-asg-dyndns/dns/cloudflare_dns/mocks"
	"github.com/cloudflare/cloudflare-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCloudflareDomainToZoneResolve(t *testing.T) {
	domain := "test.test1.example.com"
	expected := "example.com"

	var cfp cloudflareProvider
	cfp.detectZone(domain)

	assert.Equal(t, expected, cfp.zone, "Zone should be equal.")
}

func TestNewCloudflareProvider(t *testing.T) {
	ctx := context.TODO()
	expectedToken := "token123ABC"
	expectedName := "TOKEN_PARAM_NAME"
	domain := "www.test.example.com"
	expectedZone := "example.com"

	ctrl := gomock.NewController(t)
	m := awsClient_mocks.NewMockSSMAPI(ctrl)

	fakeSSMResponse := ssm.GetParameterOutput{
		Parameter: &types.Parameter{
			Name:  &expectedName,
			Value: &expectedToken,
		},
	}
	m.EXPECT().GetParameter(gomock.Any(), gomock.Any()).Return(&fakeSSMResponse, nil)

	provider, err := NewCloudflareProvider(ctx, m, domain)

	assert.Equalf(t, expectedToken, provider.token, "Token %s not equal to expected value %s", provider.token, expectedToken)
	assert.Equalf(t, expectedZone, provider.zone, "Zone %s not equal to expected value %s", provider.zone, expectedZone)
	assert.Nilf(t, err, "Error %s not null", err)
}

func TestDnsEntryAddIp(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	// aws_mock := awsClient_mocks.NewMockSSMAPI(ctrl)
	cf_mock := cf_mocks.NewMockCLOUDFLAREAPI(ctrl)

	zone := "example.com"
	domain := "test.example.com"
	ip := "190.191.10.12"

	zone_id_value := "adadad1231313"
	var zone_id_err error = nil

	cf_records := []cloudflare.DNSRecord{
		// cloudflare.DNSRecord{
		// 	Type:       "A",
		// 	Name:       domain,
		// 	Content:    ip,
		// 	CreatedOn:  time.Time{},
		// 	ModifiedOn: time.Time{},
		// 	ID:         "123123",
		// 	TTL:        60,
		// 	ZoneID:     zone_id_value,
		// 	ZoneName:   zone,
		// },
	}
	var cf_recors_err error = nil

	dns_record_response := cloudflare.DNSRecordResponse{}
	var create_dns_record_error error = nil

	cf_mock.EXPECT().ZoneIDByName(gomock.Eq(zone)).Return(zone_id_value, zone_id_err)
	cf_mock.EXPECT().DNSRecords(gomock.Any(), gomock.Eq(zone_id_value), gomock.Any()).
		Return(cf_records, cf_recors_err)
	cf_mock.EXPECT().CreateDNSRecord(gomock.Any(), gomock.Eq(zone_id_value), gomock.Any()).
		Return(&dns_record_response, create_dns_record_error)

	cf := cloudflareProvider{
		token: "test_token_123",
		api:   cf_mock,
		zone:  zone,
	}

	err := cf.DnsEntryAddIp(ctx, domain, &ip)
	assert.Nilf(t, err, "Expected no error, but got %s", err)
}

func TestDnsEntryRemoveIp(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	cf_mock := cf_mocks.NewMockCLOUDFLAREAPI(ctrl)

	zone := "example.com"
	domain := "test.example.com"
	ip := "190.191.10.12"

	zone_id_value := "adadad1231313"
	var zone_id_err error = nil

	cf_records := []cloudflare.DNSRecord{
		cloudflare.DNSRecord{
			Type:       "A",
			Name:       domain,
			Content:    ip,
			CreatedOn:  time.Time{},
			ModifiedOn: time.Time{},
			ID:         "123123",
			TTL:        60,
			ZoneID:     zone_id_value,
			ZoneName:   zone,
		},
	}
	var cf_recors_err error = nil

	var delete_dns_record_error error = nil

	cf_mock.EXPECT().ZoneIDByName(gomock.Eq(zone)).Return(zone_id_value, zone_id_err)
	cf_mock.EXPECT().DNSRecords(gomock.Any(), gomock.Eq(zone_id_value), gomock.Any()).
		Return(cf_records, cf_recors_err)
	cf_mock.EXPECT().DeleteDNSRecord(gomock.Any(), gomock.Eq(zone_id_value), gomock.Any()).
		Return(delete_dns_record_error)

	cf := cloudflareProvider{
		token: "test_token_123",
		api:   cf_mock,
		zone:  zone,
	}

	err := cf.DnsEntryRemoveIp(ctx, domain, &ip)
	assert.Nilf(t, err, "Expected no error, but got %s", err)
}
