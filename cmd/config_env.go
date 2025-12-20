// Copyright © 2025 sstreichan
// This file is part of github.com/sstreichan/ztdns.

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// loadEnvConfig reads configuration from environment variables (ZTDNS_ prefix)
// and returns values used by the application. Missing required variables return an error.
func loadEnvConfig() error {
	// Load .env file if present (values are used only when the same var is not set in the real environment)
	if err := loadDotEnv(); err != nil {
		// Non-fatal: if .env missing, continue; other errors report
		if !os.IsNotExist(err) {
			return fmt.Errorf("loading .env: %w", err)
		}
	}

	// Required
	if os.Getenv("ZT_API") == "" {
		return fmt.Errorf("missing required env ZT_API")
	}
	if os.Getenv("ZT_URL") == "" {
		return fmt.Errorf("missing required env ZT_URL")
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

	domainsEnv := os.Getenv("ZT_DOMAIN")
	networksEnv := os.Getenv("ZT_NETWORK")
	if domainsEnv == "" || networksEnv == "" {
		return fmt.Errorf("missing required envs: ZT_DOMAIN and ZT_NETWORK must be set")
	}
	// split and trim
	splitAndTrim := func(s string) []string {
		parts := strings.Split(s, ",")
		out := []string{}
		for _, x := range parts {
			x = strings.TrimSpace(x)
			if x != "" {
				out = append(out, x)
			}
		}
		return out
	}
	domains := splitAndTrim(domainsEnv)
	networks := splitAndTrim(networksEnv)
	if len(domains) != len(networks) {
		return fmt.Errorf("ZT_DOMAIN and ZT_NETWORK must have the same number of comma-separated entries")
	}
	for i, domain := range domains {
		nid := networks[i]
		key := "ZT_NETWORK_" + strings.ToUpper(strings.ReplaceAll(domain, "-", "_"))
		os.Setenv(key, nid)
	}
	return nil
}

// loadDotEnv loads a .env file at repo root and sets variables only when they are not already present in the environment.
func loadDotEnv() error {
	f, err := os.Open(".env")
	if err != nil {
		return err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			// skip malformed lines
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), "\"')")
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
	if err := s.Err(); err != nil {
		return err
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
