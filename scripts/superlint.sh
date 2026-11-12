#!/usr/bin/env bash

# Script must be inside git repo to work (which should always be true, since the script is provided in the github repository)
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# The abuslute repository path should be found, no matter where inside the repository the file is placed
repo_dir="$(cd "${script_dir}" && git rev-parse --show-toplevel 2>/dev/null)" || {
  echo "[FATAL]: Script is not inside a valid git repository." >&2
  exit 1
}

docker run -e LOG_LEVEL=DEBUG \
 -e RUN_LOCAL=true \
 -e DEFAULT_BRANCH=main \
 -v "${repo_dir}"://tmp/lint \
 --rm ghcr.io/super-linter/super-linter:latest@sha256:acc7cc64c938cb76c7239f2a862674cdc0078776f280bdd6bb6802c693e0cf3b
docker_status=$?

if [ "$docker_status" -ne 0 ]; then
  if [ "$docker_status" -eq 125 ]; then
    echo "[INFO]: The docker daemon responded with an error (Docker desktop may not be running)"
  elif [ "$docker_status" -eq 126 ]; then
    echo "[INFO]: A docker command within superlint.sh could not be invoked"
  elif [ "$docker_status" -eq 127 ]; then
    echo "[INFO]: A docker command within superlint.sh could not be found or does not exist"
  fi
  echo "[FATAL]: Docker exited code $docker_status"
fi
echo "[INFO]: Linter succesfully Completed"
exit $docker_status