test:
	@./go.test.sh
.PHONY: test

coverage:
	@./go.coverage.sh
.PHONY: coverage

test_fast:
	go test ./...

tidy:
	go mod tidy

build:
	GOAMD64=v3 go build -o ./_bin/ ./cmd/*

install:
	go install ./cmd/...

deploy-worker:
	cat hosts.workers.txt | xargs -I HOST scp _bin/gha-worker HOST:~

deploy-restic:
	cat hosts.workers.txt | xargs -I HOST scp /usr/local/bin/restic HOST:~

deploy-controller:
	cat hosts.controller.txt | xargs -I HOST scp _bin/gha-controller HOST:~

deploy: build install deploy-worker deploy-controller
