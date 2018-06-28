#!/bin/bash

set -exuo pipefail

go build -o iaam5-app-db main.go
docker build --tag symflower/iaam5:app-db .
