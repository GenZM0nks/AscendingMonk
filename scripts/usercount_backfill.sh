#!/usr/bin/env bash
set -euo pipefail
umask 177

backup_file=""

cleanup () {
    if [[ -n "$backup_file"]]; then
        rm -f -- "$backup_file"
    fi
}

skip () {
    echo "Warning: $1 for $backup_filename from $backup_date"
    cleanup
}

trap cleanup EXIT
trap cleanup ERR

USERCOUNT_LOG=/var/log/ascendingmonk/usercount.csv
DB_BLOB_DIR=/tmp
TIMESTAMP=$(date -u +"%FT%T.%7NZ")

if [[ -f "$USERCOUNT_LOG" ]]; then
    echo "Error: A log file already exists at: $USERCOUNT_LOG" >&2
    exit 1
fi
echo "Info: Creating log file for user count at: $USERCOUNT_LOG"
touch "$USERCOUNT_LOG"

printf "%s\n" "logged_at,source,backup_date,user_count" >"$USERCOUNT_LOG"

if ! azure_response=$(az storage blob list --container-name database-backups --prefix daily/ --auth-mode login --query '[].name' --output tsv); then
    echo "Error: Failed to retrive azure storage names" >&2
    exit 1
fi

mapfile -t azure_data <<< "$azure_response"

backup_files=$(sort -k2,2n <<< "$azure_data")

for backup_blob in "${backup_files[@]}";
do 
    if [[ -z "$backup_blob" ]]; then
        echo "Warning: backup name empty"
        continue
    fi
    backup_date=$(cut -d/ -f2 <<< $backup_blob)
    backup_filename=$(cut -d/ -f3 <<<$backup_blob)
    backup_file="$DB_BLOB_DIR/$backup_filename"
    if ! az storage blob download --account-name examplebackupstorage --container-name database-backups --name "$backup_blob" --file "$backup_file" --auth-mode login --overwrite false --only-show-errors --output none; then
        echo "Warning: Failed to retrive azure storage blob: $backup_filename"
        continue
    fi
    if [[ "$(sqlite3 "$backup_file" "PRAGMA integrity_check";)" != "ok" ]]; then
        skip "integrity failed"
        continue
    fi
    echo "Info: integrity verified for $backup_filename from $backup_date"
    if ! query_response=$(sqlite3 "$backup_file" "SELECT COUNT(*) FROM users"); then
        skip "failed to query user count"
        continue
    fi
    if [[ ! "$query_response" =~ ^[0-9]+$ ]]; then
        skip "Unexpected query response ( $query_response )"
        continue
    fi
    printf "%s,%s,%s,%s\n" "$TIMESTAMP" "Azure blob storage" "$backup_date" "$query_response" >>"$USERCOUNT_LOG"
    cleanup
done