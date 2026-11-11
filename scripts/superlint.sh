#!/usr/bin/env bash

docker run \
  -e LOG_LEVEL=DEBUG \
  -e RUN_LOCAL=true \
  -v ../:/tmp/lint \
  --rm \
  ghcr.io/super-linter/super-linter:latest@sha-acc7cc64c938cb76c7239f2a862674cdc0078776f280bdd6bb6802c693e0cf3b # acc7cc64c938cb76c7239f2a862674cdc0078776f280bdd6bb6802c693e0cf3b

docker_status=$?

echo $docker_status

if [ "$docker_status" -eq 125 ]; then
  echo "[INFO]: "
if [ "$docker_status" -eq 125 ]; then
  echo "[INFO]: The docker daemon responded with an error (Docker desktop may not be running)"
elif [ "$docker_status" -eq 126 ]; then
  echo "[INFO]: A docker command within superlint.sh could not be invoked"
elif [ "$docker_status" -eq 127 ]; then
  echo "[INFO]: A docker command within superlint.sh could not be found or does not exist"
fi
