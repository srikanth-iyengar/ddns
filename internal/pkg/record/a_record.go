package record

import (
	"encoding/binary"

	"github.com/srikanth-iyengar/ddns/internal/pkg/dns"
)

type ARecord struct {
	dns.ResourcePreamble
	Ip uint32
}

func (aRecord ARecord) WireFormat() []byte {
	buffer := aRecord.Query.WireFormat()
	buffer = append(buffer, 0, 0, 0, 0)

	binary.BigEndian.PutUint32(buffer[len(buffer)-4:], aRecord.Ttl)
	buffer = append(buffer, 0, 4)
	buffer = append(buffer, aRecord.Data()...)
	return buffer
}

func (aRecord ARecord) Data() []byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, aRecord.Ip)
	return buffer
}

func (aRecord ARecord) Preamble() dns.ResourcePreamble {
	return aRecord.ResourcePreamble
}
