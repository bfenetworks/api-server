# Copyright (c) 2019 The BFE Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# init project path
HOMEDIR := $(shell pwd)
OUTDIR  := $(HOMEDIR)/output

# -----------------------------------------------------------------------------
# Container image build variables (used by docker/docker-push targets)
# -----------------------------------------------------------------------------

# Image name and version
BFE_IMAGE_NAME ?= api-server
BFE_VERSION ?= $(shell sed -nE 's/^const Version = "([^"]+)"/\1/p' version/version.go)
VERSION_TAG := v$(BFE_VERSION)

# Registry and build settings
REGISTRY ?=
PLATFORMS ?= linux/amd64,linux/arm64
BUILDER_NAME ?= api-server-builder
NO_CACHE ?= false
DOCKER_BUILD_ARGS ?=

# Derived tags
IMAGE_LOCAL := $(BFE_IMAGE_NAME):$(VERSION_TAG)
IMAGE_LATEST_LOCAL := $(BFE_IMAGE_NAME):latest
IMAGE_REMOTE := $(if $(REGISTRY),$(REGISTRY)/$(BFE_IMAGE_NAME):$(VERSION_TAG),)
IMAGE_LATEST_REMOTE := $(if $(REGISTRY),$(REGISTRY)/$(BFE_IMAGE_NAME):latest,)

DOCKER_NO_CACHE_FLAG := $(if $(filter true,$(NO_CACHE)),--no-cache,)

# init command params
GO_BIN ?= go
GO      := $(GO_BIN)
GOPATH  := $(shell $(GO) env GOPATH)
GOMOD   := $(GO) mod
GOBUILD := $(GO) build
GOTEST  := $(GO) test -gcflags="-N -l"
GOPKGS  := $$($(GO) list ./...| grep -vE "vendor")
LICENSEEYE   := license-eye

# test cover files
COVPROF := $(HOMEDIR)/covprof.out  # coverage profile
COVFUNC := $(HOMEDIR)/covfunc.txt  # coverage profile information for each function
COVHTML := $(HOMEDIR)/covhtml.html # HTML representation of coverage profile

# make, make all
all: prepare compile package

#make prepare, download dependencies
prepare: gomod

gomod:
	$(GOMOD) download

#make compile
compile: build

build: prepare
	$(GOBUILD) -o $(HOMEDIR)/api-server

# make test, test your code
test: prepare test-case
test-case:
	$(GOTEST) -v -cover $(GOPKGS)

# make package
package: package-bin
package-bin:
	mkdir -p 			$(OUTDIR)
	cp -rf  conf 		$(OUTDIR)/
	cp -rf  static 		$(OUTDIR)/
	cp -rf  docs 		$(OUTDIR)/
	cp -rf  db_ddl.sql 	$(OUTDIR)/
	cp -rf  *.md 		$(OUTDIR)/
	mv api-server  		$(OUTDIR)/

# make docker
docker:
	docker build $(if $(filter 1,$(NO_CACHE)),--no-cache,) $(DOCKER_BUILD_ARGS) \
		-t $(IMAGE_LOCAL) \
		-t $(IMAGE_LATEST_LOCAL) \
		.

# make docker-push (multi-arch)
docker-push:
	@test -n "$(REGISTRY)" || (echo "REGISTRY is required, e.g. REGISTRY=ghcr.io/your-org" && exit 1)
	@docker buildx inspect $(BUILDER_NAME) >/dev/null 2>&1 || docker buildx create --name $(BUILDER_NAME) --driver docker-container --use
	@docker buildx use $(BUILDER_NAME)
	@docker buildx inspect --bootstrap >/dev/null
	docker buildx build $(if $(filter 1,$(NO_CACHE)),--no-cache,) $(DOCKER_BUILD_ARGS) \
		--platform $(PLATFORMS) \
		--push \
		-t $(IMAGE_REMOTE) \
		-t $(IMAGE_LATEST_REMOTE) \
		.

# make license-eye-install
license-eye-install:
	$(GO) install github.com/apache/skywalking-eyes/cmd/license-eye@latest

# make license-check, check code file's license declaration
license-check: license-eye-install
	$(LICENSEEYE) header check

# make license-fix, fix code file's license declaration
license-fix: license-eye-install
	$(LICENSEEYE) header fix

# make clean
clean:
	$(GO) clean
	rm -rf $(OUTDIR)
	rm -rf $(HOMEDIR)/api-server
	rm -rf $(GOPATH)/pkg/darwin_amd64

# avoid filename conflict and speed up build 
.PHONY: all prepare compile test package clean build docker docker-push
