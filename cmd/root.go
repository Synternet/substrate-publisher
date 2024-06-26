package cmd

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"github.com/synternet/data-layer-sdk/pkg/network"
	"github.com/synternet/data-layer-sdk/pkg/options"
)

var (
	flagVerbose       *bool
	flagNatsUrls      *string
	flagUserCreds     *string
	flagNkey          *string
	flagJWT           *string
	flagTLSClientCert *string
	flagTLSKey        *string
	flagCACert        *string
	flagPrefixName    *string

	natsConnection *nats.Conn
)

var rootCmd = &cobra.Command{
	Use:   "substrate-publisher",
	Short: "",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetFlags(0)
		var err error

		err = network.SetDefault("testnet")
		if err != nil {
			panic(err)
		}

		natsConnection, err = options.MakeNats("Substrate Publisher", *flagNatsUrls, *flagUserCreds, *flagNkey, *flagJWT, *flagCACert, *flagTLSClientCert, *flagTLSKey)
		if err != nil {
			panic(err)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if natsConnection == nil {
			return
		}
		natsConnection.Close()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	const (
		PUBLISHER_IDENTITY = "PUBLISHER_IDENTITY"
		STREAM_PREFIX      = "STREAM_PREFIX"
	)
	setDefault(STREAM_PREFIX, "synternet")

	flagNatsUrls = rootCmd.PersistentFlags().StringP("nats-urls", "n", os.Getenv("NATS_URLS"), "NATS server URLs (separated by comma)")
	flagUserCreds = rootCmd.PersistentFlags().StringP("nats-creds", "c", os.Getenv("NATS_CREDS"), "NATS User Credentials File (combined JWT and NKey file) ")
	flagJWT = rootCmd.PersistentFlags().StringP("nats-jwt", "w", os.Getenv("NATS_JWT"), "NATS user JWT")
	flagNkey = rootCmd.PersistentFlags().StringP("nats-nkey", "k", os.Getenv("NATS_NKEY"), "NATS user NKey")
	flagTLSKey = rootCmd.PersistentFlags().StringP("client-key", "", "", "NATS Private key file for client certificate")
	flagTLSClientCert = rootCmd.PersistentFlags().StringP("client-cert", "", "", "NATS TLS client certificate file")
	flagCACert = rootCmd.PersistentFlags().StringP("ca-cert", "", "", "NATS CA certificate file")
	flagPrefixName = rootCmd.PersistentFlags().StringP("stream-prefix", "", os.Getenv(STREAM_PREFIX), "NATS topic prefix name as in {prefix}.substrate")
	flagVerbose = rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
}
