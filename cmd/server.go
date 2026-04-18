/*
Copyright © 2026 Srikanth Iyengar srikanth.iyengar@srikanthk.in
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/spf13/cobra"
	"github.com/srikanth-iyengar/ddns/config"
	"github.com/srikanth-iyengar/ddns/pkg/api"
	"github.com/srikanth-iyengar/ddns/pkg/handler"
	v1 "github.com/srikanth-iyengar/ddns/proto/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	grpcConfig               config.ServerConfig
	dnsConfig                config.ServerConfig
	shouldRegisterReflection bool
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start dns grpc server",
	Long: `Usage: 
		dns server start
	`,
	Run: func(cmd *cobra.Command, args []string) {
		lis, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", grpcConfig.Hostname, grpcConfig.Port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		dnsService := api.DnsResourceServer{}
		grpcServer := grpc.NewServer()
		v1.RegisterDnsServiceServer(grpcServer, &dnsService)
		if shouldRegisterReflection {
			log.Printf("Registering reflection...")
			reflection.Register(grpcServer)
		}
		serverGroup, ctx := errgroup.WithContext(context.Background())

		serverGroup.Go(func() error {
			return grpcServer.Serve(lis)
		})

		serverGroup.Go(func() error {
			return handler.ServeDns(fmt.Sprintf(":%d", dnsConfig.Port))
		})

		<-ctx.Done()
		if err := serverGroup.Wait(); err != nil {
			log.Fatalf("failed: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVar(&grpcConfig.Hostname, "grpc-hostname", "127.0.0.1", "Bind hostname for grpc server")
	serverCmd.Flags().Int16Var(&grpcConfig.Port, "grpc-port", 8000, "Port configuration for the grpc server")
	serverCmd.Flags().BoolVar(&shouldRegisterReflection, "reflection", true, "Register reflection service in grpc server")
	serverCmd.Flags().StringVar(&dnsConfig.Hostname, "dns-hostname", "127.0.0.1", "Bind hostname for dns server")
	serverCmd.Flags().Int16Var(&dnsConfig.Port, "dns-port", 53, "Port configuration for the dns server")
}
