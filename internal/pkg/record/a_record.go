package record

import "github.com/srikanth-iyengar/ddns/internal/pkg/dns"

type ARecord struct {
	dns.ResourcePreamble
	Ip uint32
}

func (aRecord *ARecord) WireFormat() []byte {
	return aRecord.ResourcePreamble.WireFormat()
}
