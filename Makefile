SHELL := /usr/bin/env bash

.DEFAULT_GOAL := help

## --------------------------------------
## Help
## --------------------------------------

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

test:
	go test -timeout 30s ./... -v

build:
	CGO_ENABLED=0 go build -o bin/vmware-work-sample main.go

image:
	docker build -t ghcr.io/jtcressy/vmware-work-sample:dev -f Dockerfile .

deploy-default:
	kubectl apply -k config/overlays/default

deploy-google-pinger:
	kubectl apply -k config/overlays/google-pinger

deploy: deploy-default

deploy-local: minikube image
	kubectl apply -k config/overlays/local-minikube

minikube:
	minikube status || make minikube-start

minikube-start:
	minikube start

minikube-stop:
	minikube stop

template:
	kubectl apply -k config/overlays/default --dry-run=client -o yaml