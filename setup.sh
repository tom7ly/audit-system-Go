#!/bin/bash

set -e

start_postgresql_container() {
  echo "Checking if PostgreSQL container is already running..."
  if [ $(docker ps -q -f name=postgres-audit-system) ]; then
    echo "PostgreSQL container already running."
  else
    echo "Starting PostgreSQL container..."
    docker pull postgres:latest
    docker run --name postgres-audit-system -e POSTGRES_USER=pq -e POSTGRES_PASSWORD=pq -e POSTGRES_DB=audit -p 5432:5432 -d postgres:latest

    echo "Waiting for PostgreSQL to be ready..."
    for i in {1..30}; do
      if docker exec postgres-audit-system pg_isready -U pq -d audit > /dev/null 2>&1; then
        echo "PostgreSQL is ready."
        break
      else
        echo "Waiting for PostgreSQL to become ready..."
        sleep 10
      fi
    done

    if [ $i -eq 30 ]; then
      echo "PostgreSQL failed to become ready in time."
      exit 1
    fi

    docker exec postgres-audit-system psql -U pq -d postgres -c "CREATE DATABASE test_audit;"
    echo "PostgreSQL container started and configured."
  fi
}

install_go() {
  echo "Installing Go 1.22.4..."
  wget https://dl.google.com/go/go1.22.4.linux-amd64.tar.gz
  sudo tar -C /usr/local -xzf go1.22.4.linux-amd64.tar.gz
  rm go1.22.4.linux-amd64.tar.gz
  echo "Go 1.22.4 installed."
}

run_go_application() {
  echo "Setting up Go application..."
  cd .

    go mod tidy

    go run -mod=mod entgo.io/ent/cmd/ent generate ./ent/schema

    go run cmd/main.go
}

main() {
    start_postgresql_container
    install_go
    run_go_application
}

main
