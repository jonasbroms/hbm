#!/usr/bin/env bash

ROOTDIR=$(dirname $0)/../../..
cd $(dirname $0)

if [ -d "build" ]; then
	rm -rf build
fi
mkdir -p build

cp ${ROOTDIR}/scripts/packages/hbm.service build/
cp ${ROOTDIR}/scripts/packages/hbm.socket build/
cp ${ROOTDIR}/bin/hbm build/

go run ${ROOTDIR}/gen/man/genman.go
cp -r /tmp/hbm/man build/

go run ${ROOTDIR}/gen/shellcompletion/genshellcompletion.go
cp -r /tmp/hbm/shellcompletion build/
