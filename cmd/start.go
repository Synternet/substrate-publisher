package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/synternet/data-layer-sdk/pkg/service"
	"github.com/synternet/substrate-publisher/internal/substrate"
)

var (
	flagTendermintAPI *string
	flagRPCAPI        *string
	flagGRPCAPI       *string
	flagPublisherName *string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		publisher := substrate.New(
			service.WithContext(ctx),
			service.WithName(*flagPublisherName),
			service.WithPrefix(*flagPrefixName),
			service.WithNats(natsConnection),
			service.WithUserCreds(*flagUserCreds),
			service.WithNKeySeed(*flagNkey),
			service.WithVerbose(*flagVerbose),
			substrate.WithRPCAPI(*flagRPCAPI),
		)

		if publisher == nil {
			return
		}

		pubCtx := publisher.Start()
		defer publisher.Close()

		select {
		case <-ctx.Done():
			log.Println("Shutdown")
		case <-pubCtx.Done():
			log.Println("Publisher stopped with cause: ", context.Cause(pubCtx).Error())
		}
	},
}

func setDefault(field string, value string) {
	if os.Getenv(field) == "" {
		os.Setenv(field, value)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)

	const (
		SUBSTRATE_RPC  = "SUBSTRATE_RPC"
		PUBLISHER_NAME = "PUBLISHER_NAME"
	)

	setDefault(SUBSTRATE_RPC, "wss://rpc.polkadot.io")
	setDefault(PUBLISHER_NAME, "polkadot")

	flagPublisherName = startCmd.Flags().StringP("name", "", os.Getenv(PUBLISHER_NAME), "NATS subject name as in {prefix}.{name}.>")
	flagRPCAPI = startCmd.Flags().StringP("app-api", "a", os.Getenv(SUBSTRATE_RPC), "Full address to the Applications RPC")
}
