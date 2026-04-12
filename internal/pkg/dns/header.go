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
	ID         // request ID
	QR         // Query(0) | Response(1) bit
	OpCode     // 0 = Standard Query, 1 = Inverse Query, 2 = Status.
	AA         // Authoritative Answer - 1 if the server is an authority.
	TC         // TrunCation - 1 if the message was truncated.
	RD         // Recursion Desired - 1 if recursive query needed.
	RA         // Recursion Available - 1 if server supports recursion.
	Z          // Reserved for future use (set to 0).
	RCode      // Response Code - 0=No error, 3=Name Error (NXDOMAIN).
	QueryCount // Number of entries in the Question Section.
	AnsCount   // Number of resource records in the Answer Section.
	NSCount    // Number of name server records in Authority Section.
	ARCount    // Number of resource records in Additional Section.
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

	buffer = append(buffer, byte(h.ID>>8), byte(h.ID&0x00FF))

	var flag1 uint8

	if h.QR {
		flag1 = 0x80
	}

	flag1 |= uint8(h.OpCode) << 3

	if h.AA {
		flag1 |= 0x04
	}

	if h.TC {
		flag1 |= 0x02
	}

	if h.RD {
		flag1 |= 0x01
	}

	buffer = append(buffer, byte(flag1))

	// RA 			Z			Rcode
	// 1 				3			4
	var flag3 = uint8(h.RCode)

	flag3 |= uint8(h.Z << 4)

	if h.RA {
		flag3 |= 0x8
	}

	buffer = append(buffer, flag3)

	buffer = append(buffer, byte(h.QueryCount>>8), byte(h.QueryCount&0x00FF))
	buffer = append(buffer, byte(h.AnsCount>>8), byte(h.AnsCount&0x00FF))
	buffer = append(buffer, byte(h.NSCount>>8), byte(h.NSCount&0x00FF))
	buffer = append(buffer, byte(h.ARCount>>8), byte(h.ARCount&0x00FF))

	return buffer
}
