package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// loadEnvConfig reads configuration from environment variables (ZTDNS_ prefix)
// and returns values used by the application. Missing required variables return an error.
func loadEnvConfig() error {
	// Required
	if os.Getenv("ZTDNS_ZT_API") == "" {
		return fmt.Errorf("missing required env ZTDNS_ZT_API")
	}
	if os.Getenv("ZTDNS_ZT_URL") == "" {
		return fmt.Errorf("missing required env ZTDNS_ZT_URL")
	}
	if os.Getenv("ZTDNS_NETWORKS") == "" {
		return fmt.Errorf("missing required env ZTDNS_NETWORKS (format: domain=networkid,domain2=networkid2)")
	}

	// Optional with defaults
	if os.Getenv("ZTDNS_SUFFIX") == "" {
		os.Setenv("ZTDNS_SUFFIX", "zt")
	}
	if os.Getenv("ZTDNS_PORT") == "" {
		os.Setenv("ZTDNS_PORT", "53")
	}
	if os.Getenv("ZTDNS_DBREFRESH") == "" {
		os.Setenv("ZTDNS_DBREFRESH", "30")
	}

	// Convert NETWORKS into viper-compatible map via environment variables
	// expected format: domain=networkid,domain2=networkid2
	networks := os.Getenv("ZTDNS_NETWORKS")
	pairs := strings.Split(networks, ",")
	for _, p := range pairs {
		if p == "" {
			continue
		}
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid NETWORKS pair: %s", p)
		}
		domain := strings.TrimSpace(parts[0])
		nid := strings.TrimSpace(parts[1])
		// expose as ZT_NETWORK_<DOMAIN>=networkid to be accessible
		key := "ZT_NETWORK_" + strings.ToUpper(strings.ReplaceAll(domain, "-", "_"))
		os.Setenv(key, nid)
	}

	return nil
}

// helper to read int env
func envInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}
