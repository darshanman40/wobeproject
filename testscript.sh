#!/bin/sh

test=`go test ./... | grep "FAIL"`

if [[ -z test ]]; then
	echo "Test failed"
	echo $test
	exit 1
fi



