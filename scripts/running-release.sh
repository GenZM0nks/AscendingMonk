#!/usr/bin/env bash
set -euo pipefail

DEPLOY_DIR=/opt/ascendingmonk-deploy
COMPOSE_FILE="$DEPLOY_DIR/compose.prod.yaml"
SHA_PATTERN='^[0-9a-f]{40}$'

if [[ $# -ne 1 ]]; then
  echo "Usage: running-release.sh <40-character Git commit SHA>" >&2
  exit 2
fi

IMAGE_SHA="$1"

if [[ ! "$IMAGE_SHA" =~ $SHA_PATTERN ]]; then
  echo "Error: release must be a 40-character lowercase hexadecimal SHA." >&2
  exit 2
fi

export IMAGE_SHA

running_count=0

for service in ascendingmonk nginx; do
  set +e
  container_id="$(
    docker compose \
      --project-directory "$DEPLOY_DIR" \
      --env-file /dev/null \
      -f "$COMPOSE_FILE" \
      ps -q "$service" 2>/dev/null
  )"
  compose_status=$?
  set -e

  if ((compose_status != 0)); then
    echo "Error: could not inspect the running $service container." >&2
    exit 1
  fi

  [[ -n "$container_id" ]] || continue

  set +e
  actual_image="$(
    docker inspect \
      --format '{{.Config.Image}}' \
      "$container_id" 2>/dev/null
  )"
  inspect_status=$?
  set -e

  if ((inspect_status != 0)); then
    echo "Error: could not inspect image for the running $service container." >&2
    exit 1
  fi

  if [[ "$service" == "ascendingmonk" ]]; then
    expected_image="ghcr.io/genzm0nks/ascendingmonk:sha-$IMAGE_SHA"
  else
    expected_image="ghcr.io/genzm0nks/ascendingmonk-nginx:sha-$IMAGE_SHA"
  fi

  if [[ "$actual_image" == "$expected_image" ]]; then
    running_count=$((running_count + 1))
  fi
done

if ((running_count == 2)); then
  exit 3
fi

exit 0
