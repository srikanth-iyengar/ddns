package agent

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/srikanth-iyengar/ddns/internal/pkg/dns"
	v1 "github.com/srikanth-iyengar/ddns/proto/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	InterfaceName string
	Hostname      string
	GrpcHost      string
	GrpcPort      int16
}

func getIpv4(device string) ([][]byte, error) {
	interfaces, err := net.Interfaces()

	if err != nil {
		return nil, err
	}

	var result [][]byte
	for _, dev := range interfaces {
		addrs, _ := dev.Addrs()

		if dev.Name != device {
			continue
		}

		for _, addr := range addrs {
			switch address := addr.(type) {
			case *net.IPNet:
				{
					if ipv4 := address.IP.To4(); ipv4 != nil {
						result = append(result, ipv4)
					}
				}
			}

		}
	}

	return result, nil
}

func syncDns(cfg *Config, ctx context.Context) error {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cfg.GrpcHost, cfg.GrpcPort), opts...)

	if err != nil {
		return err
	}

	defer conn.Close()

	addrs, err := getIpv4(cfg.InterfaceName)

	if err != nil {
		return err
	}

	for _, ipv4 := range addrs {
		ipv4Num := uint32(ipv4[0])<<24 + uint32(ipv4[1])<<16 + uint32(ipv4[2])<<8 + uint32(ipv4[3])

		ddnsRequest := v1.UpsertDnsRequest{
			Preamble: &v1.Preamble{
				Qname:      strings.Split(cfg.Hostname, "."),
				QueryType:  dns.A,
				QueryClass: dns.A,
				Ttl:        0xFF,
				Length:     0x00,
			},
			Data: &v1.UpsertDnsRequest_A{
				A: &v1.ARecData{
					Ip: uint32(ipv4Num),
				},
			},
		}

		ddnsClient := v1.NewDnsServiceClient(conn)
		_, err = ddnsClient.UpsertDns(ctx, &ddnsRequest)
	}

	return err
}

func WatchInterface(ctx context.Context, cfg *Config) {
	logger, _ := zap.NewProduction()
	ticker := time.NewTicker(time.Duration(time.Second * 5))
	fmt.Println("Starting")
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			err := syncDns(cfg, ctx)
			logger.Error("Error occured while doing dns sync", zap.Error(err))
		}
	}
}
