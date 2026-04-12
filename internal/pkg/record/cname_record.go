package record

import "github.com/srikanth-iyengar/ddns/internal/pkg/dns"

type CnameRecord struct {
	dns.ResourcePreamble
	LableSequence []string
}

func (cnameRecord CnameRecord) Data() []byte {
	// TODO: implement this correctly
	return make([]byte, 10)
}

func (cnameRecord CnameRecord) Preamble() dns.ResourcePreamble {
	return cnameRecord.ResourcePreamble
}

func (cnameRecord CnameRecord) WireFormat() []byte {
	return cnameRecord.ResourcePreamble.WireFormat()
}
