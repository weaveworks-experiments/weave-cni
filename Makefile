BUILD:=./build
SRCS:=$(shell find . -path $(BUILD) -prune -o -path ./vendor -prune -o -name '*.go' -print)
DIR:=$(BUILD)/src/github.com/weaveworks/weave/cni/weave-ipam

.PHONY: all
all: ensure-deps weave-ipam

.PHONY: ensure-deps
ensure-deps:
	@git submodule update --init

weave-ipam: $(SRCS)
	mkdir -p $(DIR)
	cp $(SRCS) $(DIR)/
	cp -R ./vendor/* $(BUILD)/src/
	GOPATH=$(PWD)/$(BUILD) go install $(DIR)/
	cp ./build/bin/weave-ipam $@
