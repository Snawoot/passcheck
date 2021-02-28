#!/usr/bin/env bash

PASS_DIR_DEFAULT_CONST=~/.password-store
PASS_DIR_DEFAULT="${PASSWORD_STORE_DIR:-$PASS_DIR_DEFAULT_CONST}"
PASS_DIR="${PASS_DIR:-$PASS_DIR_DEFAULT}"

command_bin=$(command -v passcheck 2>/dev/null) && bin="$command_bin"
[[ -x ~/go/bin/passcheck ]] && bin=~/go/bin/passcheck
PASSCHECK=${PASSCHECK:-$bin}

while read -r fn ; do
	n="${fn#$PASS_DIR}"
	n="${n%.gpg}"
	p="$(gpg -d "$fn" 2>/dev/null)" && echo "\"${n//\"/\"\"}\",\"${p//\"/\"\"}\""
done < <(find "$PASS_DIR" -type f -name \*.gpg) | "$PASSCHECK" >/dev/null
