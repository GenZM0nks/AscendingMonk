#!/usr/bin/env bash
set -euo pipefail

DEPLOY_DIR=/opt/ascendingmonk-deploy
DEPLOY_SCRIPT="$DEPLOY_DIR/deploy.sh"
LOCK_FILE=/run/lock/ascendingmonk-deploy.lock
DEPLOY_LOG=/var/log/ascendingmonk/deployments.log
SHA_PATTERN='^[0-9a-f]{40}$'

if [[ $# -ne 1 ]]; then
  echo "Usage: deploy-with-rollback.sh <40-character Git commit SHA>" >&2
  exit 2
fi

IMAGE_SHA="$1"

if [[ ! "$IMAGE_SHA" =~ $SHA_PATTERN ]]; then
  echo "Error: release must be a 40-character lowercase hexadecimal SHA." >&2
  exit 2
fi

if [[ ! -x "$DEPLOY_SCRIPT" ]]; then
  echo "Error: deployment script is not executable: $DEPLOY_SCRIPT" >&2
  exit 1
fi

if [[ -e "$DEPLOY_LOG" && ! -f "$DEPLOY_LOG" ]]; then
  echo "Error: deployment log is not a regular file: $DEPLOY_LOG" >&2
  exit 1
fi

mkdir -p "$(dirname "$DEPLOY_LOG")"
touch "$DEPLOY_LOG"

exec 9>"$LOCK_FILE"
if ! flock -n 9; then
  echo "Error: another deployment, rollback, or comment operation is already running." >&2
  exit 1
fi

timestamp() {
  date -u '+%Y-%m-%dT%H:%M:%SZ'
}

append_event() {
  local event="$1"
  local sha="$2"
  local status="$3"

  printf '%s\t%s\t%s\t%s\n' \
    "$event" \
    "$(timestamp)" \
    "$sha" \
    "$status" >> "$DEPLOY_LOG"
}

previous_successful_sha() {
  awk -F '\t' '
    $1 == "D" && NF >= 4 { sha = $4 }
    END {
      if (sha != "") {
        print sha
      }
    }
  ' "$DEPLOY_LOG"
}

run_deploy() {
  local sha="$1"

  ASCENDINGMONK_DEPLOY_LOCK_HELD=1 "$DEPLOY_SCRIPT" "$sha"
}

PREVIOUS_SHA="$(previous_successful_sha)"
append_event "A" "$IMAGE_SHA" "STARTED"

echo "Deploying SHA $IMAGE_SHA."
if run_deploy "$IMAGE_SHA"; then
  append_event "S" "$IMAGE_SHA" "DEPLOYMENT_SUCCEEDED"
  exit 0
fi

append_event "F" "$IMAGE_SHA" "DEPLOYMENT_FAILED"
echo "Deployment failed for SHA $IMAGE_SHA." >&2

if [[ -z "$PREVIOUS_SHA" ]]; then
  echo "Error: no previous successful deployment is available for rollback." >&2
  exit 1
fi

echo "Rolling back to previous successful SHA $PREVIOUS_SHA."
append_event "R" "$PREVIOUS_SHA" "ROLLBACK_STARTED"

if run_deploy "$PREVIOUS_SHA"; then
  append_event "R" "$PREVIOUS_SHA" "ROLLBACK_SUCCEEDED"
  echo "Rollback succeeded for SHA $PREVIOUS_SHA." >&2
  exit 1
fi

append_event "R" "$PREVIOUS_SHA" "ROLLBACK_FAILED"
echo "Error: rollback failed for SHA $PREVIOUS_SHA." >&2
exit 1
