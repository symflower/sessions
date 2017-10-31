#!/bin/bash

set -e

for service in add sub mul div calc; do
	cd $HOME/symflower/src/socra/cmd/$service
	go build -v -o $service socra/cmd/$service
	docker build --tag socra-$service:latest .
	kubectl delete --ignore-not-found service $service
	kubectl delete --ignore-not-found deployment $service
	kubectl apply -f deployment.yml
	kubectl apply -f service.yml
done
