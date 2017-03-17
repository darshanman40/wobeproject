#!/bin/bash

# go get ./...


test=`go test ./...`
test_fail=`echo $test | grep "FAIL"`
if [[ ! -z $test_fail ]]; then
	echo "Test failed"
	echo $test
	exit 1
fi

exit 0
