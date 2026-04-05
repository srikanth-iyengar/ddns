package record

import "github.com/srikanth-iyengar/ddns/internal/pkg/dns"

type ARecord struct {
	dns.ResourcePreamble
	Ip uint32
}

func (aRecord *ARecord) WireFormat() []byte {
	return aRecord.ResourcePreamble.WireFormat()
}

func (aRecord ARecord) Data() []byte {
	// TODO: implement this correctly
	return make([]byte, 10)
}

func (aRecord ARecord) Preamble() dns.ResourcePreamble {
	return aRecord.ResourcePreamble
}
