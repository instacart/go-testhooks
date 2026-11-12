#!/bin/bash
set -e
cd $(dirname $0)

unset GOROOT

TMP_DIR=`mktemp -d`
function cleanup {      
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT

rungobuild() {
    local BINTEST=$1
    local CONFIG=$2
    TAGS=""
    if [[ "$CONFIG" == "notest" ]]; then
        TAGS="-tags notest"
    fi
    mkdir -p $TMP_DIR/$BINTEST
    CMD="go -C $(pwd)/$BINTEST build -o $TMP_DIR/$BINTEST/go-build-bin-$CONFIG.bin $TAGS -gcflags '-m -S' ."
    if ! eval $CMD 2> $TMP_DIR/$BINTEST/go-build-logs-$CONFIG.txt; then
        echo "ERROR: Failed to run $CMD" >&2
        cat $TMP_DIR/$BINTEST/go-build-logs-$CONFIG.txt >&2
        exit 1
    fi
}

for BINTEST in errhook valuehook; do
    mkdir -p $TMP_DIR/$BINTEST
    rungobuild $BINTEST default
    rungobuild $BINTEST notest
done

OUT_PARENT=$1
if [ "$OUT_PARENT" != "" ]; then
    # $1 is the parent directory; only its "out" subdirectory is ours to
    # replace, so unrelated files in the parent are left alone.
    OUT_DIR="$OUT_PARENT/out"
    mkdir -p -- "$OUT_PARENT"
    rm -rf -- "$OUT_DIR" || true
    cp -r -- "$TMP_DIR" "$OUT_DIR"
    rm -rf -- "$OUT_DIR/scripts"
    echo "Copied all logs and binaries to $OUT_DIR"
fi

python3 __main__.py "$TMP_DIR"

echo "ALL CHECKS PASS!"
