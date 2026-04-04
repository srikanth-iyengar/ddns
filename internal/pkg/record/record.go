package record

import "github.com/srikanth-iyengar/ddns/internal/pkg/dns"

type DnsRecord interface {
	Data() []byte
	Preamble() dns.ResourcePreamble
}
