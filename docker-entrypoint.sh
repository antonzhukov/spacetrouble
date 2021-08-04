#!/bin/sh
set -e

if [ "$1" = 'spacetrouble' ]; then
    exec spacetrouble \
    -pg.host="${POSTGRES_HOST}" \
    -pg.port="${POSTGRES_PORT}" \
    -pg.user="${POSTGRES_USER}" \
    -pg.password="${POSTGRES_PASSWORD}" \
    -pg.db="${POSTGRES_DB}"
fi

exec "$@"
