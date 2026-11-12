#!/usr/bin/env bash
set -euo pipefail

DEPLOY_DIR=/opt/ascendingmonk-deploy
COMPOSE_FILE="$DEPLOY_DIR/compose.prod.yaml"
LOCK_FILE=/run/lock/ascendingmonk-deploy.lock
DEPLOY_LOG=/var/log/ascendingmonk/deployments.log

if [[ $# -ne 1 ]]; then
  echo "Usage: deploy.sh <40-character Git commit SHA>" >&2
  exit 2
fi

IMAGE_SHA="$1"

if [[ ! "$IMAGE_SHA" =~ ^[0-9a-f]{40}$ ]]; then
  echo "Error: release must be a 40-character lowercase hexadecimal SHA." >&2
  exit 2
fi

export IMAGE_SHA

if [[ "${ASCENDINGMONK_DEPLOY_LOCK_HELD:-0}" != "1" ]]; then
  exec 9>"$LOCK_FILE"
  if ! flock -n 9; then
    echo "Error: another deployment, rollback, or comment operation is already running." >&2
    exit 1
  fi
fi

record_deployment() {
  local deployment_id previous_id

  if [[ -e "$DEPLOY_LOG" && ! -f "$DEPLOY_LOG" ]]; then
    echo "Error: deployment log is not a regular file: $DEPLOY_LOG" >&2
    exit 1
  fi

  mkdir -p "$(dirname "$DEPLOY_LOG")"
  touch "$DEPLOY_LOG"

  deployment_id="$(
    awk -F '\t' '
      $1 == "D" && $2 ~ /^[0-9]+$/ && $2 > highest { highest = $2 }
      END { print highest + 1 }
    ' "$DEPLOY_LOG"
  )"

  printf 'D\t%s\t%s\t%s\n' \
    "$deployment_id" \
    "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
    "$IMAGE_SHA" >> "$DEPLOY_LOG"

  previous_id="$(
    awk -F '\t' -v wanted="$IMAGE_SHA" -v current="$deployment_id" \
      '$1 == "D" && $4 == wanted && $2 != current { print $2; exit }' \
      "$DEPLOY_LOG"
  )"

  if [[ -n "$previous_id" ]]; then
    printf 'I\t%s\t%s\tSHA previously deployed as DEPLOY_ID %s\n' \
      "$deployment_id" \
      "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
      "$previous_id" >> "$DEPLOY_LOG"
    echo "Informational note added: SHA was previously deployed as DEPLOY_ID $previous_id."
  fi

  echo "Recorded successful deployment as DEPLOY_ID $deployment_id."
}

echo "Validating release: $IMAGE_SHA"

docker compose \
  --project-directory "$DEPLOY_DIR" \
  --env-file /dev/null \
  -f "$COMPOSE_FILE" \
  config --quiet

echo "Compose validation passed."

CUTOVER_MARKER="$DEPLOY_DIR/docker-cutover-complete"

if [[ ! -f "$CUTOVER_MARKER" ]]; then
  echo "Error: Docker cutover has not been completed." >&2
  echo "Routine deployment is blocked until the supervised cutover." >&2
  exit 3
fi

echo "Docker cutover marker found."
echo "Pulling images for release: $IMAGE_SHA"

docker compose \
  --project-directory "$DEPLOY_DIR" \
  --env-file /dev/null \
  -f "$COMPOSE_FILE" \
  pull

echo "Images pulled successfully."
echo "Updating Docker Compose services..."

docker compose \
  --project-directory "$DEPLOY_DIR" \
  --env-file /dev/null \
  -f "$COMPOSE_FILE" \
  up -d --no-build --remove-orphans

echo "Docker Compose services updated."
echo "Checking Docker Compose services..."

for service in ascendingmonk nginx; do
  container_id="$(
    docker compose \
      --project-directory "$DEPLOY_DIR" \
      --env-file /dev/null \
      -f "$COMPOSE_FILE" \
      ps -q "$service"
  )"

  if [[ -z "$container_id" ]]; then
    echo "Error: no running container found for $service." >&2
    exit 1
  fi

  container_state="$(
    docker inspect \
      --format '{{.State.Status}}' \
      "$container_id"
  )"

  if [[ "$container_state" != "running" ]]; then
    echo "Error: $service is $container_state, not running." >&2
    exit 1
  fi

  echo "$service is running."
done

echo "Checking deployed image references..."

for service in ascendingmonk nginx; do
  container_id="$(
    docker compose \
      --project-directory "$DEPLOY_DIR" \
      --env-file /dev/null \
      -f "$COMPOSE_FILE" \
      ps -q "$service"
  )"

  actual_image="$(
    docker inspect \
      --format '{{.Config.Image}}' \
      "$container_id"
  )"

  case "$service" in
    ascendingmonk)
      expected_image="ghcr.io/genzm0nks/ascendingmonk:sha-$IMAGE_SHA"
      ;;
    nginx)
      expected_image="ghcr.io/genzm0nks/ascendingmonk-nginx:sha-$IMAGE_SHA"
      ;;
    *)
      echo "Error: Unexpected service: $service" >&2
      exit 1
      ;;
  esac

  if [[ "$actual_image" != "$expected_image" ]]; then
    echo "Error: $service is using an unexpected image." >&2
    echo "Expected: $expected_image" >&2
    echo "Actual:   $actual_image" >&2
    exit 1
  fi

  echo "$service uses the expected release image."
done

echo "Checking application health..."

HEALTH_URL="http://127.0.0.1:8080/"
HEALTH_ATTEMPTS=12
HEALTH_DELAY=5

for ((attempt = 1; attempt <= HEALTH_ATTEMPTS; attempt++)); do
  if curl --fail --silent --show-error \
    --max-time 5 \
    --output /dev/null \
    "$HEALTH_URL"; then
    echo "Health check passed."
    record_deployment
    exit 0
  fi

  echo "Health check attempt $attempt/$HEALTH_ATTEMPTS failed."

  if ((attempt < HEALTH_ATTEMPTS)); then
    sleep "$HEALTH_DELAY"
  fi
done

echo "Error: application health check failed." >&2
exit 1
