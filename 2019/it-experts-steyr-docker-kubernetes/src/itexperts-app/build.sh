#!/bin/bash

set -exuo pipefail

cd $(dirname "$0")

go build -o itexperts-app main.go
docker build --tag symflower/itexperts:app .
