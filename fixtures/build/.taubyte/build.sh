#!/bin/bash

set -x

. /utils/wasm.sh

cat  /utils/wasm.sh

mv /src/sdk /sdk
mv /src/wasm/src/* /src
ls -lh /src
ls -lh /src/wasm
ls -lh /src/wasm/src

debug_build 1 "${FILENAME}"
ret=$?
echo -n $ret > /out/ret-code
exit $ret
