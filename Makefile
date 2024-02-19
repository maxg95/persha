.PHONY: build run
build:
	go build -o=/tmp/bin/web ./cmd/web
	
run: build
	/tmp/bin/web