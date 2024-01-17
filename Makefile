PORT := 5000

default: build

build: frontend
	go build \
		-ldflags "-X main.version=$(shell git describe --tags --always || echo dev)" \
		-mod=readonly \
		-trimpath

.PHONY: frontend
frontend: node_modules
	git clean -fdx pkg/frontend/assets/
	node ci/build.mjs

node_modules:
	vault2env --key secret/jenkins/fontawesome -- npm ci --include=dev

run: frontend
	envrun -- go run . --listen=:$(PORT)
