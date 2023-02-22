BINARY_NAME := go-scrub
DOCKER_REGISTRY := #if set it should finished by /
EXPORT_RESULT := false # for CI please set EXPORT_RESULT to true

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONEY: all test build

hello:
	echo "hello"
	echo "alwasy print hello does not exist"

## build:
build: ## build your project 
	go	build	-o	${BINARY_NAME}	.
	GOARCH=amd64	GOOS=darwin	go	build	-o	${BINARY_NAME}-darwin	.
	GOARCH=amd64	GOOS=linux	go	build	-o	${BINARY_NAME}-linux	.
	GOARCH=amd64	GOOS=windows	go	build	-o	${BINARY_NAME}-windows	.

## run:
run: ## run the project
	./${BINARY_NAME}

## build_and_run:
build_and_run: ## build and run
	build run

## clean:
clean: ## clean the project	add "-" before command so if error it is ignored
	go	clean
	-rm	${BINARY_NAME}
	-rm	${BINARY_NAME}-darwin
	-rm	${BINARY_NAME}-linux	
	-rm	${BINARY_NAME}-windows

## test:
test: ## run tests
	go test ./...
ifeq ($(EXPORT_RESULT),true)
	## GO111MODULE=off go get -u github.com/jstemmer/go-junit-report
	## $(eval OUTPUT_OPTIONS = | tee /dev/tty | go-junit-report -set-exit-code > junit-report.xml)		echo "no report"
endif
	go test -v -race ./... $(OUTPUT_OPTIONS)

## Docker:
docker-build: ## Use the dockerfile to build the container
	docker build --rm --tag $(BINARY_NAME) .

docker-release: ## Release the container with tag latest and version
	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):latest
	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)
	# Push the docker images
	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):latest
	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)


## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)