#!/bin/sh
set -eu

if [ "${1:-}" = "--terminal" ]; then
	shift
	cwd=${1:-}
	shift || true
	command=${*:-}

	if [ -z "$cwd" ] || [ -z "$command" ]; then
		exit 0
	fi

	exec nvr --remote-send "<Esc>:cd $cwd<CR>:terminal LAUNCH_EDITOR=$0 $command<CR>"
fi

file=${1:-}
line=${2:-1}
column=${3:-1}

if [ -z "$file" ]; then
	exit 0
fi

exec nvr --remote-silent +"call cursor($line, $column)" "$file"
