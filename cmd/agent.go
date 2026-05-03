/*
Copyright © 2026 Srikanth Iyengar srikanth.iyengar@srikanthk.in
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/srikanth-iyengar/ddns/internal/pkg/agent"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var cfg agent.Config

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "DDNS agent to sync hostnames",
	Long: `Starts DDNS agent to push changes of local ip address to remote ddns
	`,
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewProduction()
		errGroup, ctx := errgroup.WithContext(cmd.Context())

		errGroup.Go(func() error {
			agent.WatchInterface(ctx, &cfg)
			return nil
		})

		if err := errGroup.Wait(); err != nil {
			logger.Error("Error on wait group", zap.Error(err))
		}
	},
}

func init() {

	hostname, err := os.Hostname()

	if err != nil {
		hostname = "local.vm"
	}

	rootCmd.AddCommand(agentCmd)
	agentCmd.Flags().StringVar(&cfg.GrpcHost, "grpc-host", "127.0.0.1", "Server grpc hostname")
	agentCmd.Flags().Int16Var(&cfg.GrpcPort, "grpc-port", 8000, "Server grpc port")
	agentCmd.Flags().StringVar(&cfg.Hostname, "hostname", hostname, "hostname to be updated in DDNS")
	agentCmd.Flags().StringVar(&cfg.InterfaceName, "net", "eth0", "Net interface for ip be picked from")
}
