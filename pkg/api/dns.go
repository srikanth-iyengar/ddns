package api

import (
	"context"
	"errors"
	"strings"

	"github.com/srikanth-iyengar/ddns/internal/pkg/cache"
	"github.com/srikanth-iyengar/ddns/internal/pkg/dns"
	model "github.com/srikanth-iyengar/ddns/internal/pkg/record"
	"github.com/srikanth-iyengar/ddns/proto/v1"
	"go.uber.org/zap"
)

type DnsResourceServer struct {
	v1.UnimplementedDnsServiceServer
	Logger *zap.Logger
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
	server.Logger.Info("Upsert event",
		zap.String("qname", strings.Join(req.Preamble.Qname, ".")), zap.Uint32("query_type", req.Preamble.QueryType), zap.Uint32("ttl", req.Preamble.Ttl))
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
	}

	node := cache.UpsertRecord(record)
	records := make([]*v1.DnsData, 0)
	for _, record := range *node.Records() {
		var dnsData *v1.DnsData
		switch data := record.(type) {
		case model.ARecord:
			{
				dnsData = &v1.DnsData{
					Data: &v1.DnsData_A{
						A: &v1.ARecData{
							Ip: data.Ip,
						},
					},
				}
			}
		case model.CnameRecord:
			{
				dnsData = &v1.DnsData{
					Data: &v1.DnsData_Cname{
						Cname: &v1.CnameRecData{
							Label: data.LableSequence,
						},
					},
				}
			}
		case model.NsRecord:
			// TODO: push correct data
			{
			}
		}
		records = append(records, dnsData)
	}

	return &v1.UpsertDnsResponse{
		Preamble: req.Preamble,
		Status:   "Success",
		Records:  records,
	}, nil
}

func (server *DnsResourceServer) FindRecord(ctx context.Context, req *v1.FindRecordRequest) (*v1.FindRecordResponse, error) {
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

	if result == nil {
		return nil, errors.New("record not found")
	}

	responses := make([]*v1.FindRecordResponse_Record, len(result))

	for idx, rec := range result {
		record := v1.FindRecordResponse_Record{
			Preamble: &v1.Preamble{
				Qname:      rec.Preamble().Qname,
				Length:     uint32(rec.Preamble().Length),
				Ttl:        uint32(rec.Preamble().Ttl),
				QueryType:  uint32(rec.Preamble().QueryType),
				QueryClass: uint32(rec.Preamble().QueryClass),
			},
		}

		switch data := rec.(type) {
		case model.ARecord:
			record.Data = &v1.FindRecordResponse_Record_A{
				A: &v1.ARecData{
					Ip: data.Ip,
				},
			}
		case model.CnameRecord:
			record.Data = &v1.FindRecordResponse_Record_Cname{
				Cname: &v1.CnameRecData{
					Label: data.LableSequence,
				},
			}
		case model.NsRecord:
			record.Data = &v1.FindRecordResponse_Record_Ns{
				Ns: &v1.NsRecData{},
			}
		}

		responses[idx] = &record
	}

	return &v1.FindRecordResponse{
		Count:   uint32(len(responses)),
		Records: responses,
	}, nil
}
