package record

import "github.com/srikanth-iyengar/ddns/internal/pkg/dns"

type NsRecord struct {
	dns.ResourcePreamble
	LabelSequence []string
}

func (nsRecord NsRecord) Data() []byte {
	// TODO: implement this correctly
	return make([]byte, 10)
}

func (nsRecord NsRecord) Preamble() dns.ResourcePreamble {
	return nsRecord.ResourcePreamble
}
