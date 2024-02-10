NAME=acceleratedcharging
CLANG ?= clang-14
CFLAGS := -O2 -g -Wall -Werror $(CFLAGS)
sign = $(shell sha256sum module/module.prop | sed "s/  module\/module\.prop//g")
BUILD=CGO_ENABLED=0 go build -ldflags '-w -s' -ldflags '-X main.signature=$(sign)'
all: build \
clean

build:
	$(BUILD) -o module/bin/charge-current main.go
	cd module && zip -r ../build.zip *

clean:
	rm -rf ./module/bin/charge-current