#!/usr/bin/env bash
set -euo pipefail

DEPLOY_DIR=/opt/ascendingmonk-deploy
COMPOSE_FILE="$DEPLOY_DIR/compose.prod.yaml"

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
    exit 0
  fi

  echo "Health check attempt $attempt/$HEALTH_ATTEMPTS failed."

  if ((attempt < HEALTH_ATTEMPTS)); then
    sleep "$HEALTH_DELAY"
  fi
done

echo "Error: application health check failed." >&2
exit 1
