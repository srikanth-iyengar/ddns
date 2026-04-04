package cache

import (
	"testing"

	"github.com/srikanth-iyengar/ddns/internal/pkg/dns"
)

type mockRecord struct {
	preamble dns.ResourcePreamble
	data     []byte
}

func (m *mockRecord) Preamble() dns.ResourcePreamble {
	return m.preamble
}

func (m *mockRecord) Data() []byte {
	return m.data
}

func TestUpsertRecord(t *testing.T) {
	// Create a mock record for "www.example.com" A record
	rec := &mockRecord{
		preamble: dns.ResourcePreamble{
			Query: dns.Query{
				Qname:      []string{"com", "example", "www"},
				QueryType:  1, // A record
				QueryClass: 1, // IN
			},
			Ttl:    300,
			Length: 4,
		},
		data: []byte{192, 168, 1, 1}, // 192.168.1.1
	}

	UpsertRecord(rec)

	// Find the record
	found := FindRecord(rec)
	if found == nil {
		t.Fatal("Record not found after upsert")
	}

	if found.Preamble().QueryType != 1 {
		t.Errorf("Expected QueryType 1, got %d", found.Preamble().QueryType)
	}
}

func TestFindRecord(t *testing.T) {
	// First upsert
	rec := &mockRecord{
		preamble: dns.ResourcePreamble{
			Query: dns.Query{
				Qname:      []string{"com", "example", "www"},
				QueryType:  1,
				QueryClass: 1,
			},
			Ttl:    300,
			Length: 4,
		},
		data: []byte{192, 168, 1, 1},
	}

	UpsertRecord(rec)

	// Find with same record
	found := FindRecord(rec)
	if found == nil {
		t.Fatal("Record not found")
	}

	// Find with different query type
	rec2 := &mockRecord{
		preamble: dns.ResourcePreamble{
			Query: dns.Query{
				Qname:      []string{"com", "example", "www"},
				QueryType:  2, // NS record
				QueryClass: 1,
			},
			Ttl:    300,
			Length: 4,
		},
		data: []byte{},
	}

	found2 := FindRecord(rec2)
	if found2 != nil {
		t.Error("Found record with different query type")
	}
}

func TestUpsertMultipleRecords(t *testing.T) {
	rec1 := &mockRecord{
		preamble: dns.ResourcePreamble{
			Query: dns.Query{
				Qname:      []string{"com", "example", "www"},
				QueryType:  1,
				QueryClass: 1,
			},
			Ttl:    300,
			Length: 4,
		},
		data: []byte{192, 168, 1, 1},
	}

	rec2 := &mockRecord{
		preamble: dns.ResourcePreamble{
			Query: dns.Query{
				Qname:      []string{"com", "example", "mail"},
				QueryType:  1,
				QueryClass: 1,
			},
			Ttl:    300,
			Length: 4,
		},
		data: []byte{192, 168, 1, 2},
	}

	UpsertRecord(rec1)
	UpsertRecord(rec2)

	found1 := FindRecord(rec1)
	found2 := FindRecord(rec2)

	if found1 == nil || found2 == nil {
		t.Fatal("One or more records not found")
	}
}
