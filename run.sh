#!/bin/bash

CMD=$1

function main {
    if [ "$CMD" == "migrate_up" ]; then
       scripts/migrations.sh migrate_up
    elif [ "$CMD" == "migrate_drop" ]; then
        scripts/migrations.sh migrate_drop
    elif [ "$CMD" == "local_up" ]; then
        docker-compose -f ./docker-compose.yml build --no-cache
        docker-compose -f ./docker-compose.yml up --force-recreate
    elif [ "$CMD" == "dev_up" ]; then
        docker compose -f ./dev-docker-compose.yml build --no-cache
        docker compose -f ./dev-docker-compose.yml up --force-recreate
    elif [ "$CMD" == "stop" ]; then
        docker stop gits-server
        docker stop gits-redis
        docker stop gits-postgresql
        docker rm -v gits-server
        docker rm -v gits-redis
        docker rm -v gits-postgresql
        docker volume prune --force
    fi
}

main