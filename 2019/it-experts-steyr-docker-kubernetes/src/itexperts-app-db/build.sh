#!/bin/bash

set -exuo pipefail

cd $(dirname "$0")

go build -o itexperts-app-db main.go
docker build --tag symflower/itexperts:app-db .
