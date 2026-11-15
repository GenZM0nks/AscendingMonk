# Deploying and rolling back

## Normal deployment

1. Merge the approved change into `main`.
2. Open **Actions** in GitHub.
3. Select **Deploy to Azure**.
4. Click **Run workflow**.
5. Select `main` and start the workflow.

The workflow verifies the published images, connects to the vm, and runs the
`deploy-with-rollback.sh` script

## Reading the GitHub Actions messages

Open the workflow run and expand the **Deploy release to Azure** step.

Important messages include:

- `Deploying SHA ...`
- `Health check passed.`
- `Recorded successful deployment as DEPLOY_ID ...`
- `Deployment failed for SHA ...`
- `Rolling back to previous successful SHA ...`
- `Rollback succeeded for SHA ...`
- `Rollback failed for SHA ...`

A red workflow means the requested release was not accepted. If rollback
succeeded, the previous release should be active.

## Files and history

The Makefile is in the SSH user's home directory:

```text
~/Makefile
```

The deployment scripts are installed in:

```text
/opt/ascendingmonk-deploy/
```

The deployment history is stored in:

```text
/var/log/ascendingmonk/deployments.log
```

## Log record types

The raw log is tab-separated. The first field identifies the record type:

| Letter | Meaning |
|---|---|
| `D` | Successful deployment record. This receives a `DEPLOY_ID`. |
| `A` | Deployment attempt started. |
| `S` | Requested deployment completed successfully. |
| `F` | Requested deployment failed. |
| `R` | Rollback event from script. The status identifies whether rollback started, succeeded, or failed. |
| `C` | User comment attached to a `DEPLOY_ID`. |
| `I` | Informational note, such as a SHA being deployed previously. |

Only `D` records are used when the script selects the previous successful deployment.
`A`, `S`, `F`, and `R` records describe the operation history but are not
deployment candidates.

## List deployments

After SSH-ing into the VM you arrive in the directory where the Makefile is located:

```bash
make list
```

A list entry looks like:

```text
DEPLOY_ID 2 | timestamp | sha
    COMMENT_ID 1 (ADD): Known-good rendering
```

- `DEPLOY_ID` identifies a successful deployment.
- The SHA identifies the image hash.
- `COMMENT_ID` identifies one comment belonging to the deployment above.
- `ADD` and `REPLACE` describe the comment event (comments are not deleted but replaced)
- Attempt and rollback events are also available in the raw log:

```bash
sudo cat /var/log/ascendingmonk/deployments.log
```

## Roll back to the previous successful deployment

```bash
make rollback
```

This automatically selects the previous successful deployment and uses the existing
deployment mechanism to restore it.

## Deploy a selected deployment

By deployment ID:

```bash
make deploy DEPLOY_ID=2
```

By full SHA:

```bash
make deploy DEPLOY_ID=8b578cedaeaa8524d7f9de22e67ed7fb8f2d0502
```

## Add a comment

```bash
make comment DEPLOY_ID=2 TEXT='First good rendering of nav links'
```

Quotes inside the comment are preserved:

```bash
make comment DEPLOY_ID=2 TEXT='Correct "stable" navigation'
```

## Replace one specific comment

Use the deployment and comment IDs:

```bash
sudo -n /opt/ascendingmonk-deploy/rollback.sh \
  comment replace 2 1 'Updated rendering after CSS fix'
```

This appends a replacement record; it does not erase the old comment from
the raw log.

## When the workflow has a red dot

1. Read the **Deploy release to Azure** step in GitHub Actions.
2. Look for whether rollback was attempted and whether it succeeded.
3. SSH into the VM.
4. Inspect the history:

   ```bash
   make list
   sudo cat /var/log/ascendingmonk/deployments.log
   ```

5. Check the running containers and site if needed:

   ```bash
   sudo docker ps
   curl -k -I --max-time 5 https://ascendingmonk.dk/
   ```

If rollback succeeded, the previous release should be running. If rollback
failed, identify a known-good `DEPLOY_ID` and deploy it explicitly:

```bash
make deploy DEPLOY_ID=<known-good-DEPLOY_ID>
```

Do not run `make rollback` until the history and current state have
been inspected.
