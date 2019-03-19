#!/bin/sh
set -ex
package=github.com/measurement-kit/measurement-kit/go/task
gomobile bind -v -x \
  -target android \
  -javapkg io.ooni.mk.go \
  -ldflags '-s -w' \
  -v $package
gomobile bind -v -x \
  -target ios \
  -ldflags '-s -w' \
  -prefix MKGo \
  -v $package
