#!/bin/sh
set -e
rm -rf manpages
mkdir manpages
cd cli
go run . man | gzip -c > "../manpages/gsoc2.1.gz"