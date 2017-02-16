
PACKAGE=github.com/terrywh/keytracker
VERSION=0.2.0

VENDORS=${GOPATH}/src/BurntSushi/toml ${GOPATH}/src/julienschmidt/httprouter ${GOPATH}/src/gorilla/websocket

SOURCE_ENTRY=$(wildcard main/*.go)
SOURCE_FILES=$(wildcard *.go) $(wildcard */*.go)

TARGET=bin/keytracker

.PHONY: vendor test run

all: ${TARGET}

${TARGET}: vendor ${SOURCE_FILES}
	${GOROOT}/bin/go build -ldflags "-X ${PACKAGE}/config.AppVersion=${VERSION}" -o $@ ${PACKAGE}/main

vendor:

run: vendor
# 测试状态设置 GOGC 提高内存回收频率
	GOGC=10 ${GOROOT}/bin/go run -ldflags "-X ${PACKAGE}/config.AppVersion=${VERSION} -X ${PACKAGE}/config.AppPath=/data/godocs/src/${PACKAGE}" ${SOURCE_ENTRY}

clean:
	rm -f ${TARGET}
