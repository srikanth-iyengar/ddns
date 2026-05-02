package record

import (
	"encoding/binary"

	"github.com/srikanth-iyengar/ddns/internal/pkg/dns"
)

type CnameRecord struct {
	dns.ResourcePreamble
	LableSequence []string
}

func (cnameRecord CnameRecord) Data() []byte {
	res := make([]byte, 0)

	for _, label := range cnameRecord.LableSequence {
		res = append(res, byte(len(label)))
		res = append(res, []byte(label)...)
	}

	res = append(res, 0x00)

	return res
}

func (cnameRecord CnameRecord) Preamble() dns.ResourcePreamble {
	return cnameRecord.ResourcePreamble
}

func (cnameRecord CnameRecord) WireFormat() []byte {
	buffer := cnameRecord.Query.WireFormat()

	buffer = binary.BigEndian.AppendUint32(buffer, cnameRecord.Ttl)

	data := cnameRecord.Data()

	buffer = binary.BigEndian.AppendUint16(buffer, uint16(len(data)))

	buffer = append(buffer, data...)

	return buffer
}
