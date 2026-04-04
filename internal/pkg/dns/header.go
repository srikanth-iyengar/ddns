package dns

type ID uint16
type QR bool
type OpCode uint8
type AA bool
type TC bool
type RD bool
type RA bool
type Z uint8
type RCode uint8
type QueryCount uint16
type AnsCount uint16
type NSCount uint16
type ARCount uint16

type Header struct {
	ID
	QR
	OpCode
	AA
	TC
	RD
	RA
	Z
	RCode
	QueryCount
	AnsCount
	NSCount
	ARCount
}

func NewHeader(
	id ID,
	qr QR,
	opCode OpCode,
	aa AA,
	tc TC,
	rd RD,
	ra RA,
	z Z,
	rCode RCode,
	queryCount QueryCount,
	ansCount AnsCount,
	nsCount NSCount,
	arCount ARCount,
) Header {
	return Header{
		ID:         id,
		QR:         qr,
		OpCode:     opCode,
		AA:         aa,
		TC:         tc,
		RD:         rd,
		RA:         ra,
		Z:          z,
		RCode:      rCode,
		QueryCount: queryCount,
		AnsCount:   ansCount,
		NSCount:    nsCount,
		ARCount:    arCount,
	}
}

func (h *Header) WireFormat() []byte {
	buffer := make([]byte, 0)

	buffer = append(buffer, byte(h.ID>>8), byte(h.ID&0x0F))

	var qrBit uint8

	if h.QR {
		qrBit = 0x08
	}

	// QR			OpCode
	// 1 Bit	3 Bit
	buffer = append(buffer, byte(uint8(h.OpCode&0x0E)>>1|qrBit))

	// OpCode			AA 	TC	RD
	// 1 Bit			1		1 	1
	var flag2 uint8

	flag2 |= uint8(h.OpCode) & 0x01 << 3
	if h.AA {
		flag2 |= 0x01 << 2
	}

	if h.TC {
		flag2 |= 0x01 << 1
	}

	if h.RD {
		flag2 |= 0x01
	}

	buffer = append(buffer, byte(flag2))

	return buffer
}
