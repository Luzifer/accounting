PORT := 5000

default: build

build: export CGO_ENABLED=0
build: frontend
	go build \
		-ldflags "-s -w -X main.version=$(shell git describe --tags --always || echo dev)" \
		-mod=readonly \
		-trimpath

.PHONY: frontend
frontend: node_modules
	git clean -fdx pkg/frontend/assets/
	node ci/build.mjs

node_modules:
	npm ci --include=dev

publish: frontend
	bash ci/build.sh

run: frontend
	go run . --listen=:$(PORT)

trivy:
	trivy fs . \
		--dependency-tree \
		--exit-code 1 \
		--format table \
		--ignore-unfixed \
		--quiet \
		--scanners license,misconfig,secret,vuln \
		--severity HIGH,CRITICAL
