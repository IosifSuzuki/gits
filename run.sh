#!/bin/bash

CMD=$1

function main {
    if [ "$CMD" == "migrate_up" ]; then
       scripts/migrations.sh migrate_up
    elif [ "$CMD" == "migrate_drop" ]; then
        scripts/migrations.sh migrate_drop
    elif [ "$CMD" == "up" ]; then
       docker-compose build --no-cache
       docker-compose up --force-recreate
    fi
}

main