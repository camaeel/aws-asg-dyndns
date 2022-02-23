package dns

type cloudflare struct {
	token string
}

func (c cloudflare) dnsEntryAddIp(domain string, ip *string) error {
	return nil
}

func (c cloudflare) dnsEntryRemoveIp(domain string, ip *string) error {
	return nil
}
