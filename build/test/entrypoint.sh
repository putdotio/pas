#!/bin/bash -ex

mysql=(mysql -hmysql -upas -p123)
max_attempts=30
attempt=0

until "${mysql[@]}" -e "select 1" &>/dev/null ; do
  attempt=$(( attempt + 1 ))
  if [ $attempt -ge $max_attempts ]; then
    echo "Error: MySQL did not become ready within $max_attempts attempts"
    exit 1
  fi
  echo "MySQL is not ready yet... (attempt $attempt/$max_attempts)"
  sleep 1
done
echo "MySQL is ready."

exec go test -v -race -covermode atomic -coverprofile=/coverage/covprofile ./...
