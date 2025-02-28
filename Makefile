OPENAPIDIR= ./api
OPENAPI = $(OPENAPIDIR)/location.yaml
APISPEC := $(shell cat  $(OPENAPI)| grep title: | head -1 |sed 's/ /-/g' | sed 's/^.*title:-//'| tr '[:upper:]' '[:lower:]')

PROJECT_ROOT = github.com/davidsteed/cognito

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -ldflags "-X main.Version=`git rev-parse --short HEAD`" -a -o ./build/bootstrap ./cmd/
	yarn && yarn build

.PHONY: deploy
deploy: build
	cd deploy ; yarn build ; npx cdk deploy -c zoneId=ZJP01E7QBLR9Q -c zoneName=testawsreact.com -c subdomain=www --all; cd ..

.PHONY: run
run:
	go run ./cmd/

.PHONY: tidy
tidy:
	find . -name '*.go' -ipath "*/node_modules/*" -exec rm {} \;
	go mod tidy


.PHONY: openapi-build2-go
openapi-build:
	echo $(APISPEC)
	oapi-codegen --config ./config.yaml  $(OPENAPI)> $(OPENAPIDIR)/location.gen.go
	npx --yes orval@7.6.0 --config ./orval.config.cjs
