# Building & Run
build: bin
	go build -o bin/version-to-regex

run: build 
	cd ./bin && ./version-to-regex

fmt:
	go fmt ./...

tests:
# use https://github.com/mfridman/tparse/ for prettier output
	go test -v -cover -count=1 -json ./... | tparse -all -progress

test:
	go test -v ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

# Demo commands
demo: build
	@echo "=== Testing exact version match ==="
	./bin/version-to-regex "1.2.3" 1.2.3 1.2.4 1.3.0
	@echo "\n=== Testing caret range ==="
	./bin/version-to-regex "^1.2.3" 1.2.3 1.2.5 1.3.0 2.0.0
	@echo "\n=== Testing tilde range ==="
	./bin/version-to-regex "~1.2.3" 1.2.3 1.2.5 1.3.0
	@echo "\n=== Testing wildcard ==="
	./bin/version-to-regex "1.2.*" 1.2.0 1.2.999 1.3.0

bin:
	mkdir -p bin

# Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	golangci-lint run --fix

format:
	@echo "Running golangci-lint..."
	go fmt .
