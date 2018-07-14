####################################################################################
###                                                                             ####
###                                                                             ####
###                                                                             ####
####################################################################################

SOURCEDIR          =.
SOURCES            :=$(shell find $(SOURCEDIR) -name '*.go')

####################################################################################

VERSION            = $(shell git rev-parse HEAD)
BUILD_TIME         = `date +%FT%T%z`

BINARY             = main

LD_X               = -X main.VERSION=${VERSION} -X main.BUILD_TIME=${BUILD_TIME}
LD_EXT             = -s -extldflags \"-static\"
LD_FLAGS           = -w -linkmode external ${LD_X} ${LD_EXT}

####################################################################################

.PHONY: test clean

default:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go generate
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "${LD_FLAGS}" -o ${BINARY}

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "${LD_FLAGS}" -o ${BINARY}

test:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
