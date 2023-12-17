#!/bin/bash


DEV_ENV="dev"

ENV=$1
CMD=$2

function main {
    loadAppropriateENV

    if ! which migrate >/dev/null; then
      echo "warning: migrate is not installed"
    elif [ "$CMD" == "migrate_up" ]; then
        migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${PORT_DB}/${POSTGRES_DB}?sslmode=${POSTGRES_MODE}" -verbose up
    elif [ "$CMD" == "migrate_down" ]; then
        migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${PORT_DB}/${POSTGRES_DB}?sslmode=${POSTGRES_MODE}" -verbose down
    else
        migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${PORT_DB}/${POSTGRES_DB}?sslmode=${POSTGRES_MODE}" -verbose drop
    fi
}

function loadAppropriateENV {
    if [ "$ENV" == "$DEV_ENV" ]; then
      set -a
      source ./build/dev/.env
      set +a
    else
      set -a
      source ./build/prod/.env
      set +a
    fi
}

main
