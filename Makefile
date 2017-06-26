
PACKAGE=github.com/terrywh/keytracker
VERSION=1.0.0

VENDORS=${GOPATH}/src/BurntSushi/toml ${GOPATH}/src/julienschmidt/httprouter ${GOPATH}/src/gorilla/websocket

SOURCE_ENTRY=$(wildcard main/*.go)
SOURCE_FILES=$(wildcard *.go) $(wildcard */*.go)

TARGET_LINUX=bin/keytracker
TARGET_WIN64=bin/keytracker.exe

.PHONY: get test run win32

all: ${TARGET_LINUX}

${TARGET_LINUX}: ${SOURCE_FILES}
	GOOS=linux go build -ldflags "-X ${PACKAGE}/config.AppVersion=${VERSION}" -o $@ ${PACKAGE}/main
${TARGET_WIN64}: ${SOURCE_FILES}
	GOOS=windows go build -ldflags "-X ${PACKAGE}/config.AppVersion=${VERSION}" -o $@ ${PACKAGE}/main
win64: ${TARGET_WIN64}

get:
	go get github.com/BurntSushi/toml
	go get github.com/julienschmidt/httprouter
	go get github.com/gorilla/websocket
clean:
	rm -f ${TARGET_LINUX} ${TARGET_WIN64}
