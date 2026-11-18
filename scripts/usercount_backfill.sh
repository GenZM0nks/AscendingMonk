#!/usr/bin/env bash

USERCOUNT_LOG="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/test/usercount.csv" # temp path for development

if [[ -f "$USERCOUNT_LOG" ]]; then
    echo "Error: A log file already exists at: $USERCOUNT_LOG"
    exit 1
fi
echo "Info: Creating log file for user count at: $USERCOUNT_LOG"
    touch $USERCOUNT_LOG