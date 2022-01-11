MOCK_FOLDER=${PWD}/pkg/mock
COVERAGE_FILE=$(IGNORED_FOLDER)/coverage.out
IGNORED_FOLDER=.ignore

.PHONY: tools cover-html cover clean test lint mock install builde

##
## Building
##

install: ## Download and install go mod
	@go mod download

build: export GOOS=linux
build: export GOARCH=amd64
build: ## Build App
	@go build -v -o ${BIN_FOLDER}/${APP} $(shell go list -m)/

.PHONY: install build

##
## Quality Code
##

mock: ## Mock
	MOCK_FOLDER=${PWD}/mock go generate ./...

lint: ## Lint
	@golangci-lint run

test: mock
	@mkdir -p ${IGNORED_FOLDER}
	@go test -gcflags=-l -count=1 -race -coverprofile=${COVERAGE_FILE} -covermode=atomic ./...

test-fast:  ## Test-Fast
	@mkdir -p ${IGNORED_FOLDER}
	@go test -count=1 -race -coverprofile=${COVERAGE_FILE} -covermode=atomic ./...

cover: ## Cover
	@if [ ! -e ${COVERAGE_FILE} ]; then \
		echo "Error: ${COVERAGE_FILE} doesn't exists. Please run \`make test\` then retry."; \
		exit 1; \
	fi
	@go tool cover -func=${COVERAGE_FILE}

cover-html: ## Cover html
	@if [ ! -e ${COVERAGE_FILE} ]; then \
		echo "Error: ${COVERAGE_FILE} doesn't exists. Please run \`make test\` then retry."; \
		exit 1; \
	fi
	@go tool cover -html=${COVERAGE_FILE}

clean:
	@rm -rf ${IGNORED_FOLDER}
	@rm -rf ${COVERAGE_FILE}

##
## Tooling
##

tools-lint: ## Install go lint dependencies
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

tools-test: ## Install go test dependencies
	@go install github.com/golang/mock/mockgen@latest

tools: tools-lint tools-test ## Install lint test dependencies

.PHONY: tools-lint tools-test tools