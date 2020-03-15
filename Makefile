.PHONY: build clean

BINARY=rcse
VERSION=0.1
GIT_COMMIT=$(shell git rev-parse HEAD)
BUILDDATE=$(shell date -u)
LDFLAGS=-ldflags "-X 'github.com/nateph/rcse/cmd.rcse=$(VERSION)' -X 'github.com/nateph/rcse/cmd.gitCommit=$(GIT_COMMIT)'\
 -X 'github.com/nateph/rcse/cmd.buildDate=$(BUILDDATE)'"

build:
	go build -o $(BINARY) $(LDFLAGS)
install:
	go install $(LDFLAGS)
darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY) $(LDFLAGS)
linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY) $(LDFLAGS)	
test:
	go test -race -v ./...
clean:
	rm $(BINARY)