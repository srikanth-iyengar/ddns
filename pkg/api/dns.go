package api

import (
	"context"
	"log"

	"github.com/srikanth-iyengar/ddns/internal/pkg/cache"
	"github.com/srikanth-iyengar/ddns/internal/pkg/dns"
	model "github.com/srikanth-iyengar/ddns/internal/pkg/record"
	"github.com/srikanth-iyengar/ddns/proto/v1"
)

type DnsResourceServer struct {
	v1.UnimplementedDnsServiceServer
}

func (server *DnsResourceServer) UpsertDns(ctx context.Context, req *v1.UpsertDnsRequest) (*v1.UpsertDnsResponse, error) {
	var record model.DnsRecord
	preamble := dns.ResourcePreamble{
		Query: dns.Query{
			Qname:      req.Preamble.Qname,
			QueryType:  uint16(req.Preamble.QueryType),
			QueryClass: uint16(req.Preamble.QueryClass),
		},
		Ttl:    req.Preamble.Ttl,
		Length: uint16(req.Preamble.Length),
	}
	switch data := req.GetData().(type) {
	case *v1.UpsertDnsRequest_A:
		record = model.ARecord{
			Ip:               data.A.Ip,
			ResourcePreamble: preamble,
		}
	case *v1.UpsertDnsRequest_Cname:
		record = model.CnameRecord{
			LableSequence:    data.Cname.Label,
			ResourcePreamble: preamble,
		}
	case *v1.UpsertDnsRequest_Ns:
		record = model.NsRecord{
			LabelSequence:    data.Ns.LabelSequence,
			ResourcePreamble: preamble,
		}
	default:
		return nil, nil
	}

	cache.UpsertRecord(record)
	return nil, nil
}

func (server *DnsResourceServer) FindRecord(ctx context.Context, req *v1.FindDnsRequest) (*v1.UpsertDnsResponse, error) {
	preamble := dns.ResourcePreamble{
		Query: dns.Query{
			Qname:      req.Preamble.Qname,
			QueryType:  uint16(req.Preamble.QueryType),
			QueryClass: uint16(req.Preamble.QueryClass),
		},
		Ttl:    req.Preamble.Ttl,
		Length: uint16(req.Preamble.Length),
	}

	result := cache.FindRecord(&preamble)

	log.Printf("Result: %v\n", result)

	dnsResponse := v1.UpsertDnsResponse{
		Preamble: &v1.Preamble{
			Qname:      result[0].Preamble().Qname,
			Length:     uint32(result[0].Preamble().Length),
			Ttl:        result[0].Preamble().Ttl,
			QueryType:  uint32(result[0].Preamble().QueryType),
			QueryClass: uint32(result[0].Preamble().QueryClass),
		},
		Status: "Success",
	}

	return &dnsResponse, nil
}
