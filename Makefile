# Options
SERVICE_NAME=shipa-keptn
ORG_NAME ?= shipasoftware
VERSION ?= 0.0.1

build:
	go build -ldflags '-linkmode=external' -v -o shipa-keptn

image:
	docker build . -t $(ORG_NAME)/$(SERVICE_NAME):$(VERSION)

image-push:
	docker push $(ORG_NAME)/$(SERVICE_NAME):$(VERSION)

lint:
	$(LINT) run

# Tools

LINT=$(shell which golangci-lint)

.PHONY: lint build image image-push
