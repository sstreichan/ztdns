// Copyright © 2017 uxbh
// Copyright © 2025 sstreichan
// This file is part of github.com/sstreichan/ztdns.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ztdns",
	Short: "Zerotier DNS Server",
	Long: `ztDNS is a dedicated DNS server for ZeroTier networks.
This application will serve DNS requests for the members of a ZeroTier
network for both A (IPv4) and AAAA (IPv6) requests`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().Bool("debug", false, "enable debug messages")
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))

	// Config is environment-only; no --config flag
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Configuration is environment-only. Load env-based config and expose to viper for compatibility.
	viper.SetEnvPrefix("ZTDNS")
	viper.AutomaticEnv()

	// Load our env conversion helper which validates required vars and sets defaults
	if err := loadEnvConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
		os.Exit(2)
	}

	// Map simple env keys to viper
	viper.BindEnv("debug", "ZTDNS_DEBUG")
	viper.BindEnv("interface", "ZTDNS_INTERFACE")
	viper.BindEnv("port", "ZTDNS_PORT")
	viper.BindEnv("suffix", "ZTDNS_SUFFIX")
	viper.BindEnv("DbRefresh", "ZTDNS_DBREFRESH")
	viper.BindEnv("ZT.API", "ZTDNS_ZT_API")
	viper.BindEnv("ZT.URL", "ZTDNS_ZT_URL")

	// For networks, read environment variables created in loadEnvConfig
	// and set them into Viper's string map
	networks := map[string]string{}
	for _, e := range os.Environ() {
		// look for ZT_NETWORK_<DOMAIN>=networkid
		if strings.HasPrefix(e, "ZT_NETWORK_") {
			parts := strings.SplitN(e, "=", 2)
			if len(parts) == 2 {
				domain := strings.TrimPrefix(parts[0], "ZT_NETWORK_")
				domain = strings.ToLower(strings.ReplaceAll(domain, "_", "-"))
				networks[domain] = parts[1]
			}
		}
	}
	viper.Set("Networks", networks)
}
