package dns

import "fmt"

type HeaderParser [12]byte

func (h *HeaderParser) Id() ID {
	return ID(uint16(h[0])<<8 + uint16(h[1]))
}

func (h *HeaderParser) QR() QR {
	data := h[2] >> 7

	return data == 1
}

func (h *HeaderParser) OpCode() OpCode {
	return OpCode((h[2] & 0x78) >> 3)
}

func (h *HeaderParser) AA() AA {
	return (h[2]&0x04)>>2 == 1
}

func (h *HeaderParser) TC() TC {
	return (h[2]&0x02)>>1 == 1
}

func (h *HeaderParser) RD() RD {
	return h[2]&0x01 == 1
}

func (h *HeaderParser) RA() RA {
	return (h[3]&0x80)>>7 == 1
}

func (h *HeaderParser) Z() Z {
	return Z(h[3]&0x70) >> 4
}

func (h *HeaderParser) RCode() RCode {
	return RCode(h[3] & 0x0F)
}

func (h *HeaderParser) QueryCount() QueryCount {
	return QueryCount(uint16(h[4])<<8 + uint16(h[5]))
}

func (h *HeaderParser) AnsCount() AnsCount {
	return AnsCount(uint16(h[6])<<8 + uint16(h[7]))
}

func (h *HeaderParser) NSCount() NSCount {
	return NSCount(uint16(h[8])<<8 + uint16(h[9]))
}

func (h *HeaderParser) ARCount() ARCount {
	return ARCount(uint16(h[10])<<8 + uint16(h[11]))
}

func (h *HeaderParser) String() string {
	return fmt.Sprintf(`
Id: %d,
QR: %t,
OpCode: %d,
AA: %t,
TC: %t,
RD: %t,
RA: %t,
Z: %d,
RCode: %d,
QDCount: %d,
ANCount: %d,
NSCount: %d,
ARCount: %d
	`, h.Id(),
		h.QR(),
		h.OpCode(),
		h.AA(),
		h.TC(),
		h.RD(),
		h.RA(),
		h.Z(),
		h.RCode(),
		h.QueryCount(),
		h.AnsCount(),
		h.NSCount(),
		h.ARCount())
}
