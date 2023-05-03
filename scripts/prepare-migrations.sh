#!/bin/bash

## Prepare migrations: clean old files and copy all migrations from services to migrations folder.
mkdir -p ./migrations
rm -Rvf ./migrations/*.sql
cp -Rvf ./internal/**/repository/sql/migrations/*.sql ./migrations/