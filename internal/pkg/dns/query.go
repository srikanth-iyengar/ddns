package dns

import (
	"encoding/binary"
)

type Query struct {
	Qname      []string
	QueryType  uint16
	QueryClass uint16
}

func ReadQuery(buffer []byte) (uint8, Query) {
	var query = Query{}
	query.Qname = make([]string, 0)
	var offset uint8 = 0

	for buffer[offset] != 0x00 {
		query.Qname = append(query.Qname, string(buffer[offset+1:offset+uint8(buffer[offset])+1]))
		offset += uint8(buffer[offset]) + 1
	}
	offset += 1

	query.QueryType = uint16(buffer[offset])<<8 + uint16(buffer[offset+1])
	offset += 2

	query.QueryClass = uint16(buffer[offset])<<8 + uint16(buffer[offset+1])
	offset += 2

	return offset, query
}

func (q *Query) WireFormat() []byte {
	res := make([]byte, 0)

	for _, label := range q.Qname {
		res = append(res, byte(len(label)))
		res = append(res, []byte(label)...)
	}

	res = append(res, 0x00)

	buffer := make([]byte, 2)

	binary.BigEndian.PutUint16(buffer, q.QueryType)

	res = append(res, buffer...)

	binary.BigEndian.PutUint16(buffer, q.QueryClass)

	res = append(res, buffer...)

	return res
}
