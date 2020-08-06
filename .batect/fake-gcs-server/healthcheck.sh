#!/usr/bin/env sh

set -e

HOST=${HOST:-localhost}
PORT=${PORT:-80}

curl "http://$HOST:$PORT/storage/v1/b" --fail --show-error --silent
