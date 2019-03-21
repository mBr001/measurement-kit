#!/bin/sh
set -ex
package=github.com/measurement-kit/measurement-kit/go/mkgomobile
gomobile bind -v -x \
  -target android \
  -javapkg io.ooni.mk \
  -ldflags '-s -w' \
  -v $package
gomobile bind -v -x \
  -target ios \
  -ldflags '-s -w' \
  -v $package
