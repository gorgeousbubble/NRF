# NRF Project - Makefile(Golang)
# Copyright(C) 2019-2025, Team Gorgeous Bubble, All Rights Reserved.

# Golang Commands
GO 	= go
GO-BUILD = $(GO) build
GO-CLEAN = $(GO) clean
GO-TEST = $(GO) test
GO-GET = $(GO) get

# Binary Parameters
GO-BASE = $(shell pwd)
GO-BIN = $(GO-BASE)/bin
MK-BIN = $(shell mkdir -p $(GO-BIN))

# Docker Commands
DOCKER      = docker
DOCKER-BUILD = $(DOCKER) build
DOCKER-RUN   = $(DOCKER) run

# Application
APP-NAME	= nrf
APP-DIST = nrf.tar.gz
APP-PATH = ./bin/$(APP-NAME)

# Sub Directories
SUB-DIRS = app

# Build
all: test build

build:
	$(MK-BIN)
	$(GO-BUILD) -o $(GO-BIN)

build_image:
	$(DOCKER-BUILD) -t $(APP-NAME) .

dist:	build
	tar -zcvf $(APP-DIST) $(APP-PATH)

test:
	$(GO-TEST) -v -cover -benchmem -bench .
	for dir in $(SUB-DIRS); do \
        $(MAKE) -C $$dir test; \
    done

clean:
	$(GO-CLEAN)
	rm -rf $(GO-BIN)

run:
	$(GO-BUILD) -o $(GO-BIN)
	./$(GO-BIN)

run_container:
	$(DOCKER-RUN) -it --rm --name $(APP-NAME) -p 8080:8080 -d $(APP-NAME)

deps:
	$(GO-GET) -v -t -d ./...
