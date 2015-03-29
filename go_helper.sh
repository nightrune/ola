#!/bin/sh
# Help build our go examples with our library without installing
# This needs to be generated to set the export path correctly

OLD_GOPATH=$GOPATH
export GOPATH=/home/sean/ola/golang/gocode
go build $@
rc=$?
export GOPATH=$OLD_GOPATH
exit $rc
