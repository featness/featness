# The MIT License (MIT)
# Copyright (c) 2014 globo.com <appdev@corp.globo.com>
# More info at http://opensource.org/licenses/MIT

# Several parts adapted from tsuru codebase. Original Copyright below.
# Copyright 2014 tsuru authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

define check-service
    @if [ "$(nc -z localhost $2 1>&2 2> /dev/null; echo $$?)" != "0" ]; \
    then  \
        echo "\nFATAL: Expected to find $1 running on port $2\n"; \
        exit 1; \
    fi
endef

define HG_ERROR

FATAL: you need mercurial (hg) to download tsuru dependencies.
       Check INSTALL.md for details
endef

define GIT_ERROR

FATAL: you need git to download tsuru dependencies.
       Check INSTALL.md for details
endef

define BZR_ERROR

FATAL: you need bazaar (bzr) to download tsuru dependencies.
       Check INSTALL.md for details
endef

.PHONY: all check-path get hg git bzr test

all: check-path get test

# It does not support GOPATH with multiple paths.
check-path:
ifndef GOPATH
	@echo "FATAL: you must declare GOPATH environment variable, for more"
	@echo "       details, please check INSTALL.md file and/or"
	@echo "       http://golang.org/cmd/go/#GOPATH_environment_variable"
	@exit 1
endif
ifneq ($(subst ~,$(HOME),$(GOPATH))/src/github.com/globoi/featness, $(PWD))
	@echo "FATAL: you must clone featness inside your GOPATH To do so,"
	@echo "       you can run go get github.com/globoi/featness/"
	@echo "       or clone it manually to the dir $(GOPATH)/src/github.com/globoi/featness"
	@exit 1
endif
	@echo "Great! Featness path is set correctly."
	@exit 0

get: hg git bzr godep

hg:
	$(if $(shell hg), , $(error $(HG_ERROR)))

git:
	$(if $(shell git), , $(error $(GIT_ERROR)))

bzr:
	$(if $(shell bzr), , $(error $(BZR_ERROR)))

setup: godep setup-node setup-ruby

setup-node:
	@cd dashboard && npm install .
	@cd dashboard && bower install

setup-ruby:
	@bundle

godep:
	@go get github.com/tools/godep
	@godep restore ./...
	@godep go clean ./...
	@go get github.com/jteeuwen/go-bindata/...

check-test-services:
	$(call check-service,MongoDB,3333)
	$(call check-service,Redis,4444)

_go_test:
	@go clean ./...
	@godep go test ./... -v -race

build: _build_api _build_dashboard

_build_api:
	@rm -rf ./cmd/featness-api
	@godep go build -o ./cmd/featness-api ./api-server/...
	@echo "featness-api binary up-to-date."

_build_dashboard: _build_web_app
	@rm -rf ./cmd/featness-dashboard
	@godep go build -o ./cmd/featness-dashboard ./dashboard-server/...
	@echo "featness-dashboard binary up-to-date."

_build_web_app:
	@cd dashboard && grunt build
	@mkdir -p ./dashboard-server/dashboard
	@cp -r dashboard/dist/* ./dashboard-server/dashboard
	@cd dashboard-server && go-bindata dashboard/...
	@rm -rf dashboard/dist
	@rm -rf ./dashboard-server/dashboard

test: mongo_test _go_test

run-api:
	@godep go run ./api-server/main.go --config ./api/etc/local.conf

run-dashboard:
	@cd dashboard && grunt serve

kill_redis:
	@-redis-cli -p 4444 shutdown

redis: kill_redis
	@redis-server ./redis.conf; sleep 1
	@redis-cli -p 4444 info > /dev/null

kill_mongo:
	@-ps aux | egrep -i 'mongod.+3334' | egrep -v egrep | awk '{ print $2 }' | xargs kill -9

mongo: kill_mongo
	@rm -rf /tmp/featness/mongodata && mkdir -p /tmp/featness/mongodata
	@mongod --dbpath /tmp/featness/mongodata --logpath /tmp/featness/mongolog --port 3333 --quiet &

kill_mongo_test:
	@-ps aux | egrep -i 'mongod.+3334' | egrep -v egrep | awk '{ print $2 }' | xargs kill -9

mongo_test: kill_mongo_test
	@rm -rf /tmp/featness/mongotestdata && mkdir -p /tmp/featness/mongotestdata
	@mongod --dbpath /tmp/featness/mongotestdata --logpath /tmp/featness/mongotestlog --port 3334 --quiet &

deps: mongo redis
