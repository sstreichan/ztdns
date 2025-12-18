#!/usr/bin/env bash
set -euo pipefail

# Load .env if present. Do not overwrite existing environment variables.
if [ -f ".env" ]; then
  while IFS='=' read -r key val || [ -n "$key" ]; do
    # Trim whitespace
    key="$(echo "$key" | tr -d '[:space:]')"
    # Skip empty or commented lines
    case "$key" in
      ''|#*) continue ;;
    esac
    # Remove surrounding quotes from value
    val="${val%\"}"
    val="${val#\"}"
    val="${val%\'}"
    val="${val#\'}"
    # Only set if not already in environment
    if [ -z "${!key:-}" ]; then
      export "$key"="$val"
    fi
  done < .env
fi

# Run the server (ensure the binary is built)
exec ./ztdns server
