#!bin/bash

set -exuo pipefail

while [ 1 ]
do
	tavor --format-file csvler.tavor fuzz > tmp.csv
	cat tmp.csv
	cat tmp.csv | go run main.go
done
