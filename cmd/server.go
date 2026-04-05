/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/spf13/cobra"
	"github.com/srikanth-iyengar/ddns/config"
	"github.com/srikanth-iyengar/ddns/pkg/api"
	v1 "github.com/srikanth-iyengar/ddns/proto/v1"
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
		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", grpcConfig.Hostname, grpcConfig.Port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		dnsService := api.DnsResourceServer{}
		grpcServer := grpc.NewServer()
		v1.RegisterDnsServiceServer(grpcServer, &dnsService)
		if shouldRegisterReflection {
			println("Registering reflection...")
			reflection.Register(grpcServer)
		}
		grpcServer.Serve(lis)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVar(&grpcConfig.Hostname, "grpc-hostname", "127.0.0.1", "Bind hostname for grpc server")
	serverCmd.Flags().Int16Var(&grpcConfig.Port, "grpc-port", 8000, "Port configuration for the grpc server")
	serverCmd.Flags().BoolVar(&shouldRegisterReflection, "reflection", true, "Register reflection service in grpc server")
	serverCmd.Flags().StringVar(&dnsConfig.Hostname, "dns-hostname", "127.0.0.1", "Bind hostname for dns server")
	serverCmd.Flags().Int16Var(&dnsConfig.Port, "dns-port", 8053, "Port configuration for the dns server")
}
