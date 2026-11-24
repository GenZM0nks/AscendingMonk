#!/usr/bin/env bash

# Script must be inside git repo to work (which should always be true, since the script is provided in the github repository)
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# The aboslute repository path should be found, no matter where inside the repository the file is placed
repo_dir="$(cd "${script_dir}" && git rev-parse --show-toplevel 2>/dev/null)" || {
  echo "Error: Script is not inside a valid git repository." >&2
  exit 1
}

if docker info 1>/dev/null 2>&1; then
  echo "Info: Docker deamon is running"
else
  docker_status=$?
  echo "Error: Docker daemon could no be reached. Docker exited status code $docker_status (Docker desktop must be running, for the script to succesfully run)"
  exit "$docker_status"
fi

if docker run -e LOG_LEVEL=INFO \
 -e RUN_LOCAL=true \
 --env-file "$repo_dir/.github/super-linter.env" \
 -v "${repo_dir}"://tmp/lint \
 --rm ghcr.io/super-linter/super-linter:latest@sha256:acc7cc64c938cb76c7239f2a862674cdc0078776f280bdd6bb6802c693e0cf3b; then
  echo "Info: Linter has succesfully Completed"
else 
  linter_status=$?
  echo "Error: Linter exited code $linter_status"
  exit 1  
fi

