// Copyright © 2017 uxbh
// This file is part of github.com/uxbh/ztdns.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// mkconfigCmd represents the mkconfig command
var mkconfigCmd = &cobra.Command{
	Use:   "mkconfig",
	Short: "Make a new config file",
	Long: `mkconfig (ztdns mkconfig) creates a new configuation file.
If you do not specify a filename the default is ./.ztdns.toml

Example: ztdns mkconfig [.filename.toml]`,
	Run: func(cmd *cobra.Command, args []string) {
		// mkconfig now prints example environment variables
		fmt.Println("# Example environment variables for ztdns")
		fmt.Println("export ZTDNS_SUFFIX=zt")
		fmt.Println("export ZTDNS_PORT=53")
		fmt.Println("export ZTDNS_INTERFACE=zt0")
		fmt.Println("export ZTDNS_DBREFRESH=30")
		fmt.Println("export ZTDNS_ZT_API=<<YourAPIKey>>")
		fmt.Println("export ZTDNS_ZT_URL=https://my.zerotier.com/api")
		fmt.Println("# Networks: domain=networkid,comma separated")
		fmt.Println("export ZTDNS_NETWORKS=domain=networkid,domain2=networkid2")
	},
}

func init() {
	RootCmd.AddCommand(mkconfigCmd)
}
