# How To rotate and download extended logs from the vm

Collect interaction logs daily initially, or more frequently if disk usage increases. Run only one collection at a time and avoid concurrent deployments or rotation jobs.

Use Git Bash on Windows, Bash on Linux, or run `bash` in macOS Terminal. Keep downloaded logs outside the Git repository in restricted storage. Never publish logs or private connection details.

Run each command block in order. Stop if a command or verification fails. Keep the original SSH session open throughout collection.

## 1. Connect and inspect

**Local:** Replace the placeholder with your SSH destination or configured SSH alias.

```bash
VM='YOUR_SSH_DESTINATION'
ssh "$VM"
```

**VM:**

```bash
NGINX='ascendingmonk-deploy-nginx-1'
LOG_DIR=/var/log/nginx
sudo docker ps
sudo docker exec "$NGINX" nginx -t
sudo ls -lhi "$LOG_DIR"
df -h "$LOG_DIR"
```

Before proceeding, verify that:

- The NGINX container is named `ascendingmonk-deploy-nginx-1` and its status is `Up (healthy)` (the displayed uptime may vary).
- The AscendingMonk application container is running.
- The filesystem containing `/var/log/nginx` has sufficient available disk space for compression and temporary archive copies.
- `interactions.jsonl` exists and contains data (its size is greater than zero).
- `searches.log` exists. This is the original search log and must not be renamed, truncated or deleted during collection.

If any of these checks fail, stop and investigate before proceeding with log rotation. If the interaction log is empty, there is nothing to collect.

## 2. Rotate the interaction log

**VM:**

```bash
ARCHIVE="interactions-$(date -u +%Y%m%dT%H%M%SZ).jsonl"
sudo stat -c '%n | inode=%i | bytes=%s' \
  "$LOG_DIR/interactions.jsonl" "$LOG_DIR/searches.log"

sudo bash -c '
  set -e
  cd /var/log/nginx
  test -s interactions.jsonl && test ! -L interactions.jsonl || exit 1
  test ! -e "$1" && test ! -e "$1.gz" || exit 1
  mv -- interactions.jsonl "$1"
  chmod 600 "$1"
  install -m 600 /dev/null interactions.jsonl
  docker exec "$2" nginx -s reopen
' _ "$ARCHIVE" "$NGINX"
```

If rotation fails, retain both files and investigate before retrying. Do not repeat the rename blindly.

## 3. Verify rotation

**VM:**

```bash
sudo ls -lhi "$LOG_DIR"
sudo docker top "$NGINX" -eo pid,args
```

**Identify the NGINX processes**

From the `docker top` output, identify the host process IDs (PIDs) of the NGINX master process and all worker processes.

For example, if the master PID is `12345` and the worker PID is `12346`, run:

```bash
sudo ls -l /proc/12345/fd/ /proc/12346/fd/
```

Replace the example PIDs with the actual PIDs from your output. Include every worker process if there is more than one.

**Verify open file descriptors**

Look for entries pointing to:

- `/var/log/nginx/interactions.jsonl`
- `/var/log/nginx/searches.log`

Confirm that the master and all worker processes have both active log files open and that none of them has the timestamped interaction archive open.

**Verify file inodes and permissions**

Compare the file information displayed before and after rotation. An inode identifies a file within its filesystem, even when the file is renamed.

Confirm that:

- The timestamped archive has the original interaction log's inode.
- The new `interactions.jsonl` has a different inode.
- `searches.log` retains its original inode.
- Both active logs have permissions `600`, displayed as `-rw-------`.

**Verify continued logging**

Run:

```bash
sudo wc -l "$LOG_DIR/interactions.jsonl" \
    "$LOG_DIR/searches.log"
```

Wait approximately 60 seconds and run the command again.

The interaction log should grow when eligible interactions occur. The original search log only grows when qualifying search requests occur.

If no eligible traffic occurs, wait for additional natural traffic. Do not generate artificial requests.

**Confirm that it is safe to proceed**

Only proceed with compression once the archive is no longer held open by NGINX, the active logs have the expected permissions, and continued logging has been verified.

If file descriptors cannot be inspected, inodes do not match the expected values or permissions are incorrect, stop and investigate before proceeding.

## 4. Compress and stage the archive

**VM:**

```bash
TRANSFER_DIR="$HOME/ascendingmonk-log-transfer"
mkdir -p "$TRANSFER_DIR"
chmod 700 "$TRANSFER_DIR"
(
  set -euo pipefail
  umask 077
  cd "$TRANSFER_DIR"
  test ! -e "$ARCHIVE.gz" && test ! -e "$ARCHIVE.gz.part"
  RAW_HASH=$(sudo sha256sum "$LOG_DIR/$ARCHIVE" | cut -d ' ' -f 1)
  sudo gzip -6 -c "$LOG_DIR/$ARCHIVE" > "$ARCHIVE.gz.part"
  gzip -t "$ARCHIVE.gz.part"
  test "$RAW_HASH" = "$(gzip -dc "$ARCHIVE.gz.part" | sha256sum | cut -d ' ' -f 1)"
  mv "$ARCHIVE.gz.part" "$ARCHIVE.gz"
  echo "Archive: $ARCHIVE.gz"
  echo "Raw SHA256: $RAW_HASH"
  sha256sum "$ARCHIVE.gz"
  sudo wc -l "$LOG_DIR/$ARCHIVE"
)
```

Record the filename, both hashes and record count privately. Retain the raw archive on the VM until local verification succeeds.

## 5. Download and verify

**Local — second terminal:** Replace both placeholders using the SSH destination and archive filename from the preceding steps.

```bash
VM='YOUR_SSH_DESTINATION'
ARCHIVE='interactions-REPLACE_WITH_TIMESTAMP.jsonl'
umask 077
mkdir -p "$HOME/ascendingmonk-logs/archives"
(
  set -euo pipefail
  cd "$HOME/ascendingmonk-logs/archives"
  test ! -e "$ARCHIVE.gz" && test ! -e "$ARCHIVE.gz.part"
  scp "$VM:ascendingmonk-log-transfer/$ARCHIVE.gz" "$ARCHIVE.gz.part"
  gzip -t "$ARCHIVE.gz.part"
  if command -v sha256sum >/dev/null 2>&1; then
    hash256() { sha256sum "$@"; }
  else
    hash256() { shasum -a 256 "$@"; }
  fi
  hash256 "$ARCHIVE.gz.part"
  gzip -dc "$ARCHIVE.gz.part" | hash256
  gzip -dc "$ARCHIVE.gz.part" | wc -l
)
```

Compare both hashes and the record count against the VM output. If all match, finalize the local file:

```bash
cd "$HOME/ascendingmonk-logs/archives" &&
  test ! -e "$ARCHIVE.gz" &&
  mv "$ARCHIVE.gz.part" "$ARCHIVE.gz"
```


## 6. Remove verified VM archives

**VM — original SSH session:** Only after successful local verification, confirm the raw hash still matches and the archive remains closed.

```bash
sudo sha256sum "$LOG_DIR/$ARCHIVE"
sha256sum "$TRANSFER_DIR/$ARCHIVE.gz"
```

If both hashes match the verified local values:

```bash
sudo rm -i -- "$LOG_DIR/$ARCHIVE"
rm -i -- "$TRANSFER_DIR/$ARCHIVE.gz"
sudo ls -lhi "$LOG_DIR"
df -h "$LOG_DIR"
```

Never delete either active log or use wildcard deletion. Confirm both active logs remain and continue logging.

## 7. Merge collected logs locally

Merge only verified, non-overlapping production archives. Retain the compressed originals and ensure sufficient local disk space.

```bash
(
  set -euo pipefail
  cd "$HOME/ascendingmonk-logs/archives"
  shopt -s nullglob
  files=(interactions-*.jsonl.gz)
  ((${#files[@]} > 0))
  tmp=$(mktemp ../interactions-merged.XXXXXX)
  trap 'rm -f "$tmp"' EXIT
  gzip -dc "${files[@]}" > "$tmp"
  wc -l "$tmp"
  mv "$tmp" ../interactions-merged.jsonl
)
```

Validate JSON records, duplicate request IDs and capture flags before analysis. Reconstruct approximate sessions across archive boundaries, not separately per file.

## Validation and limitations

Rotation, continued logging, gzip level 6 compression, SCP transfer, checksum verification and archive removal were successfully demonstrated on the production VM. Peak-load performance and long-term storage growth remain untested. The pre-creation of the new log with mode `600` in step 2 is a preventive refinement to the tested rotation procedure; verify it on the next collection. The revised command sequence has not been tested end-to-end on all supported operating systems.
