#!/bin/bash

scripts/migrations.sh migrate_up
docker-compose build --no-cache
docker-compose up --force-recreate
