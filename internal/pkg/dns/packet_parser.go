package dns

type PacketParser struct {
	buffer []byte
	length uint8
	HeaderParser
}

func NewPacketParser(buffer []byte, length uint8) PacketParser {
	return PacketParser{
		buffer,
		length,
		HeaderParser(buffer[:12]),
	}
}

func (msg *PacketParser) NewHeaderParser() HeaderParser {
	return HeaderParser(msg.buffer[:12])
}

func (msg *PacketParser) Queries() (uint8, []Query) {
	queries := make([]Query, 0)

	var offset uint8 = 12
	header := msg.NewHeaderParser()

	for itr := QueryCount(0); itr < header.QueryCount(); itr += 1 {
		length, query := ReadQuery(msg.buffer[offset:])
		queries = append(queries, query)
		offset += length
	}

	return offset, queries
}
