package dns

type Message struct {
	Header
	Queries []Query
}

func NewMessage(parser PacketParser) Message {
	headerParser := parser.NewHeaderParser()

	_, queries := parser.Queries()

	header := NewHeader(headerParser.Id(),
		headerParser.QR(),
		headerParser.OpCode(),
		headerParser.AA(),
		headerParser.TC(),
		headerParser.RD(),
		headerParser.RA(),
		headerParser.Z(),
		headerParser.RCode(),
		headerParser.QueryCount(),
		headerParser.AnsCount(),
		headerParser.NSCount(),
		headerParser.ARCount(),
	)

	return Message{
		Queries: queries,
		Header:  header,
	}
}
