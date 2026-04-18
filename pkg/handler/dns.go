package handler

import (
	"github.com/srikanth-iyengar/ddns/internal/pkg/cache"
	"github.com/srikanth-iyengar/ddns/internal/pkg/dns"

	"log"
	"net"
)

func ServeDns(socket string) error {

	s, err := net.ResolveUDPAddr("udp4", socket)

	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp4", s)

	if err != nil {
		return err
	}

	defer conn.Close()

	buffer := make([]byte, 512)

	for {
		n, addr, err := conn.ReadFromUDP(buffer)

		if err != nil {
			return err
		}

		parser := dns.NewPacketParser(buffer, uint8(n))
		req_header := parser.HeaderParser
		_, queries := parser.Queries()

		for _, query := range queries {
			preamble := dns.ResourcePreamble{
				Query:  query,
				Ttl:    0,
				Length: 0,
			}

			result := cache.FindRecord(&preamble)

			result_header := dns.NewHeader(
				req_header.Id(),
				dns.QR(true),
				dns.OpCode(0),
				dns.AA(false),
				dns.TC(false),
				dns.RD(req_header.RD()),
				dns.RA(false),
				dns.Z(req_header.Z()),
				dns.RCode(0),
				dns.QueryCount(1),
				dns.AnsCount(len(result)),
				dns.NSCount(0),
				dns.ARCount(0),
			)

			result_wire := make([]byte, 0)
			result_wire = append(result_wire, result_header.WireFormat()...)
			result_wire = append(result_wire, query.WireFormat()...)
			for _, rec := range result {
				result_wire = append(result_wire, rec.WireFormat()...)
			}
			_, err := conn.WriteToUDP(result_wire, addr)

			if err != nil {
				log.Printf("Error writing to sock: %v", err)
			}

		}
	}
}
