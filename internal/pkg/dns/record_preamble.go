package dns

import "encoding/binary"

type ResourcePreamble struct {
	Query
	Ttl    uint32
	Length uint16
}

func (rr *ResourcePreamble) WireFormat() []byte {
	res := make([]byte, 0)

	res = append(res, rr.Query.WireFormat()...)

	buffer := make([]byte, 2)

	binary.BigEndian.PutUint32(buffer, rr.Ttl)

	res = append(res, buffer...)

	binary.BigEndian.PutUint16(buffer, rr.Length)

	res = append(res, buffer...)

	return res
}
