#!/usr/bin/env bash
set -euo pipefail

LOG_FILE=/var/log/ascendingmonk/deployments.log
DEPLOY_SCRIPT=/opt/ascendingmonk-deploy/deploy.sh
LOCK_FILE=/run/lock/ascendingmonk-deploy.lock
SHA_PATTERN='^[0-9a-f]{40}$'

usage() {
  cat <<'EOF'
Usage:
  rollback.sh list
  rollback.sh previous
  rollback.sh rollback
  rollback.sh deploy <DEPLOY_ID|sha>
  rollback.sh comment add <DEPLOY_ID> <comment>
  rollback.sh comment replace <DEPLOY_ID> <COMMENT_ID> <comment>

The script calls the root-owned deploy.sh for deploy and rollback operations.
EOF
}

error() {
  echo "Error: $*" >&2
  exit 1
}

require_sha() {
  local sha="$1"

  if [[ ! "$sha" =~ $SHA_PATTERN ]]; then
    error "SHA must be a 40-character lowercase hexadecimal value."
  fi
}

require_deployment_id() {
  local deployment_id="$1"

  if [[ ! "$deployment_id" =~ ^[0-9]+$ ]] || ((deployment_id < 1)); then
    error "DEPLOY_ID must be a positive integer."
  fi
}

ensure_log() {
  if [[ ! -e "$LOG_FILE" ]]; then
    error "deployment log does not exist: $LOG_FILE"
  fi

  if [[ ! -f "$LOG_FILE" ]]; then
    error "deployment log is not a regular file: $LOG_FILE"
  fi
}

deployment_lines() {
  awk -F '\t' '$1 == "D" && NF >= 4 { print }' "$LOG_FILE"
}

deployment_by_id() {
  local deployment_id="$1"

  deployment_lines | awk -F '\t' -v wanted="$deployment_id" '$2 == wanted { print; exit }'
}

deployment_by_sha() {
  local sha="$1"

  deployment_lines | awk -F '\t' -v wanted="$sha" '$4 == wanted { print; exit }'
}

next_deployment_id() {
  local highest_id

  highest_id="$(
    awk -F '\t' '
      $1 == "D" && $2 ~ /^[0-9]+$/ && $2 > highest { highest = $2 }
      END { print highest + 0 }
    ' "$LOG_FILE"
  )"

  echo "$((highest_id + 1))"
}

timestamp() {
  date -u '+%Y-%m-%dT%H:%M:%SZ'
}

append_line() {
  local line="$1"
  local parent_directory

  parent_directory="$(dirname "$LOG_FILE")"
  mkdir -p "$parent_directory"
  printf '%s\n' "$line" >> "$LOG_FILE"
}

with_deployment_lock() {
  exec 9>"$LOCK_FILE"
  flock -n 9 ||
    error "another deployment, rollback, or comment operation is already running."
}

run_deploy() {
  local sha="$1"

  [[ -x "$DEPLOY_SCRIPT" ]] ||
    error "deployment script is not executable: $DEPLOY_SCRIPT"

  ASCENDINGMONK_DEPLOY_LOCK_HELD=1 "$DEPLOY_SCRIPT" "$sha"
}

validate_comment() {
  local comment="$1"

  [[ -n "$comment" ]] || error "comment cannot be empty."
  [[ "$comment" != *$'\t'* ]] || error "comment cannot contain tab characters."
  [[ "$comment" != *$'\n'* ]] || error "comment must be a single line."
  (( ${#comment} <= 500 )) ||
    error "comment is ${#comment} characters; maximum is 500. Nothing was written."
}

format_deployment() {
  local deployment_line="$1"
  local deployment_id deployment_time sha

  IFS=$'\t' read -r _ deployment_id deployment_time sha <<< "$deployment_line"
  printf 'DEPLOY_ID %s | %s | %s\n' "$deployment_id" "$deployment_time" "$sha"
}

show_comments() {
  local deployment_id="$1"

  awk -F '\t' -v wanted="$deployment_id" '
    $1 == "C" && $2 == wanted && NF >= 6 {
      printf "    COMMENT_ID %s (%s): %s\n", $3, $5, substr($0, index($0, $6))
    }
    $1 == "I" && $2 == wanted && NF >= 4 {
      printf "    info: %s\n", substr($0, index($0, $4))
    }
  ' "$LOG_FILE"
}

list_deployments() {
  ensure_log

  local count=0
  local deployment_line deployment_id

  while IFS= read -r deployment_line; do
    [[ -n "$deployment_line" ]] || continue
    format_deployment "$deployment_line"
    IFS=$'\t' read -r _ deployment_id _ _ <<< "$deployment_line"
    show_comments "$deployment_id"
    count=$((count + 1))
  done < <(deployment_lines)

  ((count > 0)) || error "deployment log contains no deployment records."
}

select_previous() {
  ensure_log

  local total
  total="$(deployment_lines | awk 'END { print NR + 0 }')"

  ((total >= 2)) ||
    error "at least two deployment records are required to select the previous deployment."

  deployment_lines | tail -n 2 | head -n 1
}

resolve_deployment() {
  local selector="$1"
  local selected

  ensure_log

  if [[ "$selector" =~ $SHA_PATTERN ]]; then
    selected="$(deployment_by_sha "$selector")"
  elif [[ "$selector" =~ ^[0-9]+$ ]]; then
    require_deployment_id "$selector"
    selected="$(deployment_by_id "$selector")"
  else
    error "selector must be a DEPLOY_ID or a 40-character lowercase hexadecimal SHA."
  fi

  [[ -n "$selected" ]] ||
    error "no deployment record matches '$selector'."

  printf '%s\n' "$selected"
}

deploy_line() {
  local selected="$1"
  local deployment_id deployment_time sha

  IFS=$'\t' read -r _ deployment_id deployment_time sha <<< "$selected"

  echo "Selected deployment:"
  format_deployment "$selected"
  echo "Deploying SHA $sha."
  run_deploy "$sha"
}

deploy_selected() {
  local selector="$1"

  deploy_line "$(resolve_deployment "$selector")"
}

next_comment_id() {
  local deployment_id="$1"
  local highest_id

  highest_id="$(
    awk -F '\t' -v wanted="$deployment_id" '
      $1 == "C" && $2 == wanted && $3 ~ /^[0-9]+$/ && $3 > highest { highest = $3 }
      END { print highest + 0 }
    ' "$LOG_FILE"
  )"

  echo "$((highest_id + 1))"
}

comment_exists() {
  local deployment_id="$1"
  local comment_id="$2"

  awk -F '\t' -v deployment="$deployment_id" -v comment="$comment_id" \
    '$1 == "C" && $2 == deployment && $3 == comment { found = 1 }
     END { exit(found ? 0 : 1) }' "$LOG_FILE"
}

add_comment() {
  local deployment_id="$1"
  local comment="$2"
  local comment_id

  require_deployment_id "$deployment_id"
  validate_comment "$comment"
  resolve_deployment "$deployment_id" >/dev/null
  comment_id="$(next_comment_id "$deployment_id")"
  append_line "C	$deployment_id	$comment_id	$(timestamp)	ADD	$comment"
  echo "Added COMMENT_ID $comment_id to DEPLOY_ID $deployment_id."
}

replace_comment() {
  local deployment_id="$1"
  local comment_id="$2"
  local comment="$3"

  require_deployment_id "$deployment_id"
  require_deployment_id "$comment_id"
  validate_comment "$comment"
  resolve_deployment "$deployment_id" >/dev/null
  comment_exists "$deployment_id" "$comment_id" ||
    error "COMMENT_ID $comment_id does not exist for DEPLOY_ID $deployment_id."
  append_line "C	$deployment_id	$comment_id	$(timestamp)	REPLACE	$comment"
  echo "Replaced COMMENT_ID $comment_id on DEPLOY_ID $deployment_id."
}

if (( $# < 1 )); then
  usage
  exit 2
fi

case "$1" in
  list)
    (( $# == 1 )) || error "list does not accept additional arguments."
    list_deployments
    ;;
  previous)
    (( $# == 1 )) || error "previous does not accept additional arguments."
    format_deployment "$(select_previous)"
    ;;
  rollback)
    (( $# == 1 )) || error "rollback does not accept additional arguments."
    with_deployment_lock
    deploy_line "$(select_previous)"
    ;;
  deploy)
    (( $# == 2 )) || error "usage: deploy <DEPLOY_ID|sha>"
    with_deployment_lock
    deploy_selected "$2"
    ;;
  comment)
    case "${2:-}" in
      add)
        (( $# == 4 )) || error "usage: comment add <DEPLOY_ID> <comment>"
        with_deployment_lock
        add_comment "$3" "$4"
        ;;
      replace)
        (( $# == 5 )) || error "usage: comment replace <DEPLOY_ID> <COMMENT_ID> <comment>"
        with_deployment_lock
        replace_comment "$3" "$4" "$5"
        ;;
      *)
        error "comment operation must be 'add' or 'replace'."
        ;;
    esac
    ;;
  -h|--help)
    usage
    ;;
  *)
    error "unknown command '$1'. Use --help for usage."
    ;;
esac
