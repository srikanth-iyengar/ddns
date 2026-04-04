package record

import "github.com/srikanth-iyengar/ddns/internal/pkg/dns"

type CnameRecord struct {
	dns.ResourcePreamble
	LableSequence []string
}
