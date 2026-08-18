#!/bin/sh

docker run -d \
    --name notes-app-db \
    -p 5432:5432 \
    -e POSTGRES_DB=notesappdb \
    -e POSTGRES_USER=admin \
    -e POSTGRES_PASSWORD=1234 \
    -v $(pwd)/init.sql:/docker-entrypoint-initdb.d/init.sql \
    -v notes-app-db-volume:/var/lib/postgresql \
    postgres
