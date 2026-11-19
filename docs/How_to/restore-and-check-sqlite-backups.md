# AscendingMonk --- Backup and Restore How-To

Azure VM, SQLite, Azure Blob Storage, Azure Managed Identity, Azure CLI

> Do not copy a restored file over the live database while the
> application is running.

## 1. Quick recovery: find an Azure backup

SSH into the AscendingMonk VM and run:

``` bash
sudo az login --identity --output none
```

List the **10 newest hourly backups on the storage**:

``` bash
sudo az storage blob list \
  --account-name ascendingmonkbackups \
  --container-name ascendingmonk-backups \
  --prefix hourly/ \
  --auth-mode login \
  --query "reverse(sort_by([].{Name:name,Modified:properties.lastModified,Bytes:properties.contentLength}, &Modified))[:10]" \
  --output table
```

Or list the **10 newest daily backups on the storage**:

``` bash
sudo az storage blob list \
  --account-name ascendingmonkbackups \
  --container-name ascendingmonk-backups \
  --prefix daily/ \
  --auth-mode login \
  --query "reverse(sort_by([].{Name:name,Modified:properties.lastModified,Bytes:properties.contentLength}, &Modified))[:10]" \
  --output table
```

Choose the backup you want. Copy its complete **Name**

## 2. Download the chosen backup safely

Create a private restore-test directory:

``` bash
sudo install -d -m 700 /var/backups/ascendingmonk/restore-test
```

Set the backup name using the value from Section 1. **Replace the
example with your chosen blob name.**

``` bash
BLOB="hourly/2026-11-19/ascendingmonk-20261119T180001Z.db"
```

Download it to a separate file:

``` bash
sudo az storage blob download \
  --account-name ascendingmonkbackups \
  --container-name ascendingmonk-backups \
  --name "$BLOB" \
  --file /var/backups/ascendingmonk/restore-test/restored.db \
  --auth-mode login \
  --overwrite true \
  --only-show-errors \
  --output none
```

Protect the downloaded database:

``` bash
sudo chmod 600 /var/backups/ascendingmonk/restore-test/restored.db
sudo ls -lh /var/backups/ascendingmonk/restore-test/restored.db
```

## 3. Health-check the downloaded database

Run SQLite's full integrity check:

``` bash
sudo sqlite3 /var/backups/ascendingmonk/restore-test/restored.db \
  "PRAGMA integrity_check;"
```

**Expected:** `ok`

Check that the expected application tables exist and has rows

``` bash
sudo sqlite3 /var/backups/ascendingmonk/restore-test/restored.db
```

```sql
.tables
SELECT * from users;
SELECT * from session_tokens;
SELECT * from pages;
```


------------------------------------------------------------------------

## 4. Check available local backups

Local snapshots are stored in:

`/var/backups/ascendingmonk/`

`sudo ls -la /var/backups/ascendingmonk/`

List the newest files:

``` bash
sudo ls -lht /var/backups/ascendingmonk/ | head -25
```

Find only scheduled SQLite snapshots:

``` bash
sudo find /var/backups/ascendingmonk/ -maxdepth 1 -type f \
  -name 'ascendingmonk-*.db' -printf '%TY-%Tm-%Td %TH:%TM  %f\n' \
  | sort -r | sed -n '1,20p'
```

Check a local snapshot **without changing it** (replace filename):

``` bash
sudo sqlite3 \
  /var/backups/ascendingmonk/ascendingmonk-20261119T180001Z.db \
  "PRAGMA integrity_check;"
```

The local backup script retains scheduled snapshots for **24 hours**.
Older manually named files may remain because they do not match the
scheduled naming pattern.

## 5. Check whether automatic backups are active

There are three systemd timers:

  ------------------------------------------------------------------------------------
  Timer                                Schedule                Purpose
  ------------------------------------ ----------------------- -----------------------
  `ascendingmonk-backup.timer`         Every 15 minutes        Creates verified local
                                                               SQLite snapshots

  `ascendingmonk-azure-upload.timer`   Every hour at `:05` UTC Uploads newest local
                                                               snapshot to `hourly/`

  `ascendingmonk-azure-daily.timer`    Daily at `00:10` UTC    Uploads newest local
                                                               snapshot to `daily/`
  ------------------------------------------------------------------------------------

Check enabled/active state:

``` bash
systemctl is-enabled ascendingmonk-backup.timer \
  ascendingmonk-azure-upload.timer \
  ascendingmonk-azure-daily.timer

systemctl is-active ascendingmonk-backup.timer \
  ascendingmonk-azure-upload.timer \
  ascendingmonk-azure-daily.timer
```

All should be **enabled** and **active**, respectively.

See the last and next scheduled runs:

``` bash
systemctl list-timers --all --no-pager --full | grep ascendingmonk
```

Inspect recent execution logs:

``` bash
sudo journalctl -u ascendingmonk-backup.service -n 30 --no-pager
sudo journalctl -u ascendingmonk-azure-upload.service -n 30 --no-pager
sudo journalctl -u ascendingmonk-azure-daily.service -n 30 --no-pager
```

Look for `Backup completed successfully`,
`Azure upload completed successfully`, and
`Daily Azure backup completed successfully`.

**Note:** A successful `Type=oneshot` service normally displays
`inactive (dead)` after it finishes. Check the timer and journal rather
than treating that status alone as a failure.

## 6. Inspect the backup services and scripts

Show systemd service and timer definitions:

``` bash
systemctl cat ascendingmonk-backup.service ascendingmonk-backup.timer
systemctl cat ascendingmonk-azure-upload.service ascendingmonk-azure-upload.timer
systemctl cat ascendingmonk-azure-daily.service ascendingmonk-azure-daily.timer
```

Read the scripts:

``` bash
sudo cat /usr/local/sbin/backup-ascendingmonk.sh
sudo cat /usr/local/sbin/upload-ascendingmonk-azure.sh
sudo cat /usr/local/sbin/upload-ascendingmonk-daily.sh
```

Check script syntax:

``` bash
sudo bash -n /usr/local/sbin/backup-ascendingmonk.sh
sudo bash -n /usr/local/sbin/upload-ascendingmonk-azure.sh
sudo bash -n /usr/local/sbin/upload-ascendingmonk-daily.sh
```

No output means no Bash syntax errors were detected. 

## 7. Browse Azure hourly and daily backups

Sign in as the VM's managed identity (works for all of us):

``` bash
sudo az login --identity --output none
```

List **all hourly backups**:

``` bash
sudo az storage blob list \
  --account-name ascendingmonkbackups \
  --container-name ascendingmonk-backups \
  --prefix hourly/ \
  --auth-mode login \
  --query "[].{Name:name,Modified:properties.lastModified,Bytes:properties.contentLength}" \
  --output table
```

List **all daily backups**:

``` bash
sudo az storage blob list \
  --account-name ascendingmonkbackups \
  --container-name ascendingmonk-backups \
  --prefix daily/ \
  --auth-mode login \
  --query "[].{Name:name,Modified:properties.lastModified,Bytes:properties.contentLength}" \
  --output table
```

List backups from a specific date, for example November 19:

``` bash
sudo az storage blob list \
  --account-name ascendingmonkbackups \
  --container-name ascendingmonk-backups \
  --prefix hourly/2026-11-19/ \
  --auth-mode login \
  --output table
```

## 8. Verified recovery evidence

On **November 19, 2026**, this backup was successfully downloaded:

`hourly/2026-11-19/ascendingmonk-20261119T180001Z.db`

Results:

-   Azure download completed successfully.
-   Downloaded size: **1,912,832 bytes** (\~1.9 MB).
-   `PRAGMA integrity_check;` returned **`ok`**.
-   Tables found: **`pages`, `session_tokens`, `users`**.
-   The tables had rows
