
PACKAGE=github.com/terrywh/keytracker
VERSION=0.3.0

VENDORS=${GOPATH}/src/BurntSushi/toml ${GOPATH}/src/julienschmidt/httprouter ${GOPATH}/src/gorilla/websocket

SOURCE_ENTRY=$(wildcard main/*.go)
SOURCE_FILES=$(wildcard *.go) $(wildcard */*.go)

TARGET=bin/keytracker

.PHONY: get test run

all: ${TARGET}

${TARGET}: ${SOURCE_FILES}
	GOOS=linux ${GOROOT}/bin/go build -ldflags "-X ${PACKAGE}/config.AppVersion=${VERSION}" -o $@ ${PACKAGE}/main

get:
	go get github.com/BurntSushi/toml
	go get github.com/julienschmidt/httprouter
	go get github.com/gorilla/websocket
clean:
	rm -f ${TARGET}
