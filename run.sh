#!/bin/bash

CMD=$1
P1=$2

function main {
    if [ "$CMD" == "dev_migrate_up" ]; then
        scripts/migrations.sh dev migrate_up
    elif [ "$CMD" == "dev_migrate_drop" ]; then
        scripts/migrations.sh dev migrate_drop
    elif [ "$CMD" == "prod_migrate_up" ]; then
        scripts/migrations.sh prod migrate_up
    elif [ "$CMD" == "prod_migrate_drop" ]; then
        scripts/migrations.sh prod migrate_drop
    elif [ "$CMD" == "generate_pass" ]; then
        go run cmd/cli/main.go pass "$P1"
    elif [ "$CMD" == "prod_up" ]; then
        docker compose -f ./build/prod/prod-docker-compose.yml build --no-cache
        docker compose -f ./build/prod/prod-docker-compose.yml up --force-recreate
    elif [ "$CMD" == "dev_up" ]; then
        docker compose -f ./build/dev/dev-docker-compose.yml build --no-cache
        docker compose -f ./build/dev/dev-docker-compose.yml up --force-recreate
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