#!/bin/bash

set -exuo pipefail

go build -o iaam5-app main.go
docker build --tag symflower/iaam5:app .
