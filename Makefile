
PACKAGE=github.com/terrywh/ntracker
VERSION=0.1.0

SOURCE_ENTRY=main.go
SOURCE_FILES=$(wildcard *.go) $(wildcard */*.go)

TARGET=bin/ntracker

.PHONY: test run

all: ${TARGET}

${TARGET}: ${SOURCE_FILES}
	${GOROOT}/bin/go build -ldflags "-X ${PACKAGE}/config.AppVersion=${VERSION}" -o $@ ${SOURCE_ENTRY}

test:

run:
# 测试状态设置 GOGC 提高内存回收频率
	GOGC=10 ${GOROOT}/bin/go run -ldflags "-X ${PACKAGE}/config.AppVersion=${VERSION} -X ${PACKAGE}/config.AppPath=/data/godocs/src/${PACKAGE}" ${SOURCE_ENTRY}

clean:
	rm -f ${TARGET}
