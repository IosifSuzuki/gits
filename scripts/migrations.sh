#!/bin/bash

set -a
source .env
set +a

CMD=$1

function main {
    if [ "$CMD" == "migrate_up" ]; then
        migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${PORT_DB}/${POSTGRES_DB}?sslmode=${POSTGRES_MODE}" -verbose up
    elif [ "$CMD" == "migrate_down" ]; then
        migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${PORT_DB}/${POSTGRES_DB}?sslmode=${POSTGRES_MODE}" -verbose down
    else
        migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${PORT_DB}/${POSTGRES_DB}?sslmode=${POSTGRES_MODE}" -verbose drop
    fi
}

main
