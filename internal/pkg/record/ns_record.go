package record

import "github.com/srikanth-iyengar/ddns/internal/pkg/dns"

type NsRecord struct {
	dns.ResourcePreamble
	LabelSequence []string
}
