#!/bin/sh
set -eu

run_sql_dir() {
  dir="$1"

  if [ ! -d "$dir" ]; then
    return 0
  fi

  for file in "$dir"/*.sql; do
    [ -e "$file" ] || continue
    echo "running $file"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f "$file"
  done
}

run_sql_dir /docker-entrypoint-initdb.d/migrations
run_sql_dir /docker-entrypoint-initdb.d/seeds
