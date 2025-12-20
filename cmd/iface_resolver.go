// Copyright © 2025 sstreichan
// This file is part of github.com/sstreichan/ztdns.

package cmd

import (
	"errors"
	"fmt"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/sstreichan/ztdns/ztapi"
)

// ResolveInterfaceFromZT tries to find the local network interface name that
// corresponds to this node's ZeroTier-assigned IPs or physical address.
//
// Parameters are explicit to keep the function testable and free of global state.
func ResolveInterfaceFromZT(apiToken, apiHost string, networks map[string]string) (string, error) {
	if apiToken == "" || apiHost == "" {
		return "", errors.New("empty API token or host")
	}
	// Collect candidate IPs and MACs from the ZT API
	candidates := map[string]struct{}{}
	candidateMACs := map[string]struct{}{}

	for _, nw := range networks {
		lst, err := ztapi.GetMemberList(apiToken, apiHost, nw)
		if err != nil {
			// Do not fail hard; log and continue with other networks
			log.Debugf("GetMemberList failed for network %s: %v", nw, err)
			continue
		}
		for _, m := range *lst {
			// collect ip assignments
			for _, a := range m.Config.IPAssignments {
				ip := net.ParseIP(strings.TrimSpace(a))
				if ip == nil {
					continue
				}
				candidates[ip.String()] = struct{}{}
			}
			// physical addr (MAC)
			mac := strings.TrimSpace(m.Config.PhysicalAddr)
			if mac != "" {
				candidateMACs[strings.ToLower(mac)] = struct{}{}
			}
		}
	}

	// List local interfaces and try to match by IP first, then MAC
	ifs, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("unable to list interfaces: %w", err)
	}

	for _, iface := range ifs {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			var ip net.IP
			switch v := a.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil {
				continue
			}
			if _, ok := candidates[ip.String()]; ok {
				log.Debugf("Matched interface %s by IP %s", iface.Name, ip.String())
				return iface.Name, nil
			}
		}
	}

	// If no IP match, try MAC match
	for _, iface := range ifs {
		mac := iface.HardwareAddr.String()
		if mac == "" {
			continue
		}
		if _, ok := candidateMACs[strings.ToLower(mac)]; ok {
			log.Debugf("Matched interface %s by MAC %s", iface.Name, mac)
			return iface.Name, nil
		}
	}

	// No match found
	log.Infof("No matching interface found from ZeroTier API")
	return "", nil
}
