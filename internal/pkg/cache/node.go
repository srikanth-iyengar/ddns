package cache

import (
	"slices"

	"github.com/srikanth-iyengar/ddns/internal/pkg/dns"
	"github.com/srikanth-iyengar/ddns/internal/pkg/record"
)

type Node struct {
	qname   string
	child   map[string]*Node
	records []record.DnsRecord
}

func (node *Node) Records() *[]record.DnsRecord {
	return &node.records
}

var root = &Node{
	qname: ".",
	child: make(map[string]*Node),
}

func (node *Node) upsertRecord(record record.DnsRecord, depth int) *Node {
	preamble := record.Preamble()

	if depth == len(preamble.Qname) {
		for _, rec := range node.records {
			if slices.Equal(rec.Data(), record.Data()) && rec.Preamble().QueryType == record.Preamble().QueryType {
				return node
			}
		}
		node.records = append(node.records, record)
		return node
	}

	levelName := preamble.Qname[depth]

	if node.child == nil {
		node.child = make(map[string]*Node)
	}

	if _, exist := node.child[levelName]; !exist {
		node.child[levelName] = &Node{qname: levelName, child: make(map[string]*Node)}
	}

	child := node.child[levelName]
	return child.upsertRecord(record, depth+1)
}

func (node *Node) findRecord(preamble *dns.ResourcePreamble, depth int) []record.DnsRecord {
	var result []record.DnsRecord
	if depth == len(preamble.Qname) {
		for _, rec := range node.records {
			if preamble.QueryType == rec.Preamble().QueryType {
				preamble.Ttl = rec.Preamble().Ttl
				result = append(result, rec)
			}
		}
		return result
	}

	levelName := preamble.Qname[depth]

	if _, exist := node.child[levelName]; exist {
		child := node.child[levelName]
		return child.findRecord(preamble, depth+1)
	} else {
		return result
	}
}

func UpsertRecord(record record.DnsRecord) *Node {
	return root.upsertRecord(record, 0)
}

func FindRecord(preamble *dns.ResourcePreamble) []record.DnsRecord {
	return root.findRecord(preamble, 0)
}
