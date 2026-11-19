#!/usr/bin/env bash
set -euo pipefail
umask 077

backup_file=""
USERCOUNT_LOG=/var/log/ascendingmonk/usercount.csv # for testing have a test folder inside script and use "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/test/usercount.csv" 
DB_BLOB_DIR="$(mktemp -d "${TMPDIR:-/tmp}/usercount.XXXXXX")" # for testning use "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/test"
TIMESTAMP=$(date -u +"%FT%T.%7NZ")

cleanup_backup () {
    echo "Info: removing temporary files..."
    if [[ -n "$backup_file" ]]; then
        echo "Info: removing temporary backup file: $backup_file"
        rm -f -- "$backup_file"
        backup_file=""
    fi
    echo "Info: removal complete"
}

cleanup_all () {
    echo "Info: starting cleanup..."
    cleanup_backup
    if [[ -n "$DB_BLOB_DIR" && -d "$DB_BLOB_DIR" ]]; then
        echo "Info: removing temporary backup directory: $DB_BLOB_DIR"
        rm -rf -- "$DB_BLOB_DIR"
    fi
    echo "Info: cleanup complete"
}

skip () {
    if [[ -n "$backup_filename" && -n "$backup_date" ]]; then
        echo "Warning: $1 for $backup_filename from $backup_date" >&2
    else
        echo "Warning: $1 for unkown backup file and date" >&2
    fi
    cleanup_backup
}

trap cleanup_all EXIT

if [[ -f "$USERCOUNT_LOG" ]]; then
    echo "Error: a log file already exists at: $USERCOUNT_LOG" >&2
    exit 1
fi
echo "Info: creating log file for user count at: $USERCOUNT_LOG"
touch "$USERCOUNT_LOG"

printf "%s\n" "logged_at,source,backup_date,user_count" >"$USERCOUNT_LOG"

if ! azure_response=$(az storage blob list --container-name database-backups --prefix daily/ --auth-mode login --query '[].name' --output tsv | sort -t/ -k2,2); then
    echo "Error: Failed to retrive azure storage names" >&2
    exit 1
fi

# For testing use the following mock data instead and manually create tree db files with those names.
# Use the following command to populate the db files: sqlite3 ascendingmonk-20261029T001005Z.db -tabs  "INSERT INTO users (username, email, password, created_at) VALUES ('user6', 'user6@stud.kea.dk', '5f4dcc3b5aa765d61d8327deb882cf99', '2026-11-17 11:47:48.826436688+01:00');"
#azure_response="daily/2026-10-27/ascendingmonk-20261027T001006Z.db
#daily/2026-10-28/ascendingmonk-20261028T001004Z.db
#daily/2026-10-29/ascendingmonk-20261029T001005Z.db"

mapfile -t azure_data <<< "$azure_response"

for backup_blob in "${azure_data[@]}";
do 
    if [[ -z "$backup_blob" ]]; then
        echo "Warning: backup name empty" >&2
        continue
    fi
    backup_date=$(cut -d/ -f2 <<<"$backup_blob")
    backup_filename=$(cut -d/ -f3 <<<"$backup_blob")
    backup_file="$DB_BLOB_DIR/$backup_filename"
    if ! az storage blob download --account-name examplebackupstorage --container-name database-backups --name "$backup_blob" --file "$backup_file" --auth-mode login --overwrite false --only-show-errors --output none; then
        skip "failed to retrive azure storage blob: $backup_filename"
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
        skip "unexpected query response ( $query_response )"
        continue
    fi
    printf "%s,%s,%s,%s\n" "$TIMESTAMP" "usercount_backfill.sh" "$backup_date" "$query_response" >>"$USERCOUNT_LOG"
    cleanup_backup
done