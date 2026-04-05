package cache

import (
	"github.com/srikanth-iyengar/ddns/internal/pkg/dns"
	"github.com/srikanth-iyengar/ddns/internal/pkg/record"
)

type Node struct {
	qname  string
	child  map[string]*Node
	record record.DnsRecord
}

var root = &Node{
	qname: ".",
	child: make(map[string]*Node),
}

func (node *Node) upsertRecord(record record.DnsRecord, depth int) *Node {
	preamble := record.Preamble()

	if depth == len(preamble.Qname) {
		node.record = record
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

func (node *Node) findRecord(preamble *dns.ResourcePreamble, depth int) record.DnsRecord {
	if depth == len(preamble.Qname) {
		if preamble.QueryType == node.record.Preamble().QueryType {
			return node.record
		}
		return nil
	}

	levelName := preamble.Qname[depth]

	if _, exist := node.child[levelName]; exist {
		child := node.child[levelName]
		return child.findRecord(preamble, depth+1)
	} else {
		return nil
	}
}

func UpsertRecord(record record.DnsRecord) *Node {
	return root.upsertRecord(record, 0)
}

func FindRecord(preamble *dns.ResourcePreamble) record.DnsRecord {
	return root.findRecord(preamble, 0)
}
