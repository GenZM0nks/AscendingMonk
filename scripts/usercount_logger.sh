#!/usr/bin/env bash
set -euo pipefail
umask 077

USERCOUNT_LOG=/var/log/ascendingmonk/usercount.csv # for testing have a test folder inside script and use "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/usercount.csv"
DB_BLOB_DIR="$(mktemp -d "${TMPDIR:-/tmp}/usercount.XXXXXX")" # for testning use "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/test"
TIMESTAMP=$(date -u +"%FT%T.%7NZ")
BACKUP_DATE="${TIMESTAMP%%T*}"
SCRIPT_DIR=/var/log/ascendingmonk/usercount-logs # for testing have a test folder inside script and use "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/usercount-logs"
SCRIPT_LOG="$SCRIPT_DIR"/"$BACKUP_DATE"_usercount_log.txt

backup_filename=""
backup_file=""

cleanup () {
    echo "Info: starting cleanup..." >>"$SCRIPT_LOG" 2>/dev/null || true
    if [[ -n "$DB_BLOB_DIR" && -d "$DB_BLOB_DIR" ]]; then
        echo "Info: removing temporary backup directory: $DB_BLOB_DIR" >>"$SCRIPT_LOG" 2>/dev/null || true
        if ! rm -rf -- "$DB_BLOB_DIR"; then
            echo "Warning: failed to remove temporary backup directory: $DB_BLOB_DIR" >>"$SCRIPT_LOG" 2>/dev/null || true
        fi
    fi
    echo "Info: cleanup complete" >>"$SCRIPT_LOG" 2>/dev/null || true
}

fail () {
    if [[ -n "$backup_filename" && -n "$BACKUP_DATE" ]]; then
        echo "Warning: $1 for $backup_filename from $BACKUP_DATE" >>"$SCRIPT_LOG"
    else
        echo "Warning: $1 for unkown backup file and date" >>"$SCRIPT_LOG"
    fi
    exit 1
}

trap cleanup EXIT

if [[ ! -f "$USERCOUNT_LOG" ]]; then
    echo "Error: a log file recording user amount does not exist at: $USERCOUNT_LOG" >&2
    exit 1
fi

if [[ ! -d "$SCRIPT_DIR" ]]; then
    echo "Error: no log directory at: $SCRIPT_DIR" >&2
    exit 1
fi

if ! touch "$SCRIPT_LOG"; then 
    echo "Error: failed to create a log file at: $SCRIPT_LOG" >&2
    exit 1
fi

echo "Info: logger Started at $TIMESTAMP" >>"$SCRIPT_LOG"

if awk -F',' -v date="$BACKUP_DATE"  '$3 == date { found=1 } END {exit !found}' "$USERCOUNT_LOG"; then
    echo "Warning: user count has already been recorded for $BACKUP_DATE" >>"$SCRIPT_LOG"
    exit 1
fi

if ! azure_response=$(az storage blob list --container-name database-backups --prefix daily/"$BACKUP_DATE"/ --auth-mode login --query '[].name' --output tsv); then
    echo "Error: Failed to retrive azure storage names" >>"$SCRIPT_LOG"
    exit 1
fi
# For testing use the following mock data instead and manually create tree db files with those names.
# Use the following command to populate the db files: sqlite3 ascendingmonk-20261029T001005Z.db -tabs  "INSERT INTO users (username, email, password, created_at) VALUES ('user6', 'user6@stud.kea.dk', '5f4dcc3b5aa765d61d8327deb882cf99', '2026-11-17 11:47:48.826436688+01:00');"
#azure_response="daily/2026-11-19/ascendingmonk-20261030T001004Z.db
#daily/2026-11-19/ascendingmonk-20261119T001006Z.db
#daily/2026-11-19/ascendingmonk-20261119T001005Z.db"


if [[ -z "$azure_response" ]]; then 
    echo "Error: no backups found for date $TIMESTAMP" >>"$SCRIPT_LOG"
    exit 1
fi

latest_backup=$(sort <<< "$azure_response" | tail -n 1)
echo "Info: using latest backup: $latest_backup" >>"$SCRIPT_LOG"

backup_filename=$(cut -d/ -f3 <<<"$latest_backup")
backup_file="$DB_BLOB_DIR/$backup_filename"
if ! az storage blob download --account-name examplebackupstorage --container-name database-backups --name "$latest_backup" --file "$backup_file" --auth-mode login --overwrite false --only-show-errors --output none; then
    fail "failed to retrive azure storage blob: $backup_filename"
fi
if [[ ! -s "$backup_file" ]]; then
    fail "downloaded backup is empty or missing"
fi
if ! integrity_response="$(sqlite3 "$backup_file" "PRAGMA integrity_check";)"; then
    fail "failed to execute integrity check"
fi
if [[ "$integrity_response" != "ok" ]]; then
    fail "backup failed integrity check"
fi
echo "Info: integrity verified for $backup_filename from $BACKUP_DATE" >>"$SCRIPT_LOG"
if ! query_response=$(sqlite3 "$backup_file" "SELECT COUNT(*) FROM users"); then
    fail "failed to query user count"
fi
if [[ ! "$query_response" =~ ^[0-9]+$ ]]; then
    fail "unexpected query response ( $query_response )"
fi
printf "%s,%s,%s,%s\n" "$TIMESTAMP" "usercount_logger.sh" "$BACKUP_DATE" "$query_response" >>"$USERCOUNT_LOG"