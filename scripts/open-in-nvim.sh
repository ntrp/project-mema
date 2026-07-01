#!/bin/sh
set -eu

file=${1:-}
line=${2:-1}
column=${3:-1}

if [ -z "$file" ]; then
	exit 0
fi

exec nvr --remote-silent +"call cursor($line, $column)" "$file"
