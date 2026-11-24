#!/usr/bin/env bash
set -euo pipefail

DEPLOY_DIR=/opt/ascendingmonk-deploy
DEPLOY_SCRIPT="$DEPLOY_DIR/deploy.sh"
RUNNING_RELEASE_CHECK="$DEPLOY_DIR/running-release.sh"
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

require_healthy_log() {
  if [[ ! -f "$DEPLOY_LOG" ]]; then
    echo "Error: deployment log is missing: $DEPLOY_LOG" >&2
    echo "Deployment aborted before images were pulled." >&2
    exit 1
  fi

  if [[ ! -r "$DEPLOY_LOG" || ! -w "$DEPLOY_LOG" ]]; then
    echo "Error: deployment log is not readable and writable: $DEPLOY_LOG" >&2
    echo "Deployment aborted before images were pulled." >&2
    exit 1
  fi

  if ! awk -F '\t' '
    $1 == "D" && $2 ~ /^[0-9]+$/ && $4 ~ /^[0-9a-f]{40}$/ { found = 1 }
    END { exit(found ? 0 : 1) }
  ' "$DEPLOY_LOG"; then
    echo "Error: deployment log has no valid successful deployment record." >&2
    echo "Deployment aborted before images were pulled." >&2
    exit 1
  fi
}

require_healthy_log

exec 9>"$LOCK_FILE"
if ! flock -n 9; then
  echo "Error: another deployment, rollback, or comment operation is already running." >&2
  exit 1
fi

if [[ ! -x "$RUNNING_RELEASE_CHECK" ]]; then
  echo "Error: running-release check is not executable: $RUNNING_RELEASE_CHECK" >&2
  exit 1
fi

set +e
"$RUNNING_RELEASE_CHECK" "$IMAGE_SHA"
running_status=$?
set -e

if ((running_status == 3)); then
  echo "Error: SHA $IMAGE_SHA is already running." >&2
  echo "Deployment aborted: no containers or deployment log entries were changed." >&2
  exit 3
fi

if ((running_status != 0)); then
  echo "Error: could not determine whether SHA $IMAGE_SHA is already running." >&2
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

try_append_event() {
  if ! append_event "$@"; then
    echo "Warning: could not write deployment event to $DEPLOY_LOG." >&2
  fi
}

previous_successful_sha() {
  awk -F '\t' '
    $1 == "D" && $2 ~ /^[0-9]+$/ && $4 ~ /^[0-9a-f]{40}$/ { sha = $4 }
    END {
      if (sha != "") {
        print sha
      }
    }
  ' "$DEPLOY_LOG"
}

run_deploy() {
  local sha="$1"
  local recovery="${2:-0}"
  local playwright_dir="/home/azureuser/AscendingMonk"
  local base_url="https://ascendingmonk.dk"

  if [[ "$recovery" == "1" ]]; then
    if ! ASCENDINGMONK_DEPLOY_LOCK_HELD=1 \
         ASCENDINGMONK_RECOVERY=1 \
         "$DEPLOY_SCRIPT" "$sha"; then
      echo "Deployment failed; skipping Playwright tests" >&2
      return 1
    fi
  else
    if ! ASCENDINGMONK_DEPLOY_LOCK_HELD=1 "$DEPLOY_SCRIPT" "$sha"; then
      echo "Deployment failed; skipping Playwright tests" >&2
      return 1
    fi
  fi

  secrets_file="/etc/ascendingmonk/prod-secrets.env"
  [[ -r "$secrets_file" ]] || {
    echo "Secrets file is missing or unreadable: $secrets_file" >&2
    return 1
  }

  source "$secrets_file"

  if ! (cd "$playwright_dir" && BASE_URL="$base_url" npx playwright test); then
    echo "Playwright tests failed - rolling back" >&2
    return 1
  fi
}

PREVIOUS_SHA="$(previous_successful_sha)"
if ! append_event "A" "$IMAGE_SHA" "STARTED"; then
  echo "Error: could not record the deployment attempt in $DEPLOY_LOG." >&2
  echo "Deployment aborted before images were pulled." >&2
  exit 1
fi

echo "Deploying SHA $IMAGE_SHA."
if run_deploy "$IMAGE_SHA"; then
  try_append_event "S" "$IMAGE_SHA" "DEPLOYMENT_SUCCEEDED"
  exit 0
fi

try_append_event "F" "$IMAGE_SHA" "DEPLOYMENT_FAILED"
echo "Deployment failed for SHA $IMAGE_SHA." >&2

if [[ -z "$PREVIOUS_SHA" ]]; then
  echo "Error: no previous successful deployment is available for rollback." >&2
  exit 1
fi

echo "Rolling back to previous successful SHA $PREVIOUS_SHA."
try_append_event "R" "$PREVIOUS_SHA" "ROLLBACK_STARTED"

if run_deploy "$PREVIOUS_SHA" 1; then
  try_append_event "R" "$PREVIOUS_SHA" "ROLLBACK_SUCCEEDED"
  echo "Rollback succeeded for SHA $PREVIOUS_SHA." >&2
  exit 1
fi

try_append_event "R" "$PREVIOUS_SHA" "ROLLBACK_FAILED"
echo "Error: rollback failed for SHA $PREVIOUS_SHA." >&2
exit 1
