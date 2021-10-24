#!/bin/sh

set -eu

: ${GOPACKAGE:=dl}

OAPI_FILE="../../doc/data-library.yaml"
OAPI_CODEGEN_VER="v1.8.3"
OAPI_CODEGEN_URL="github.com/deepmap/oapi-codegen/cmd/oapi-codegen@$OAPI_CODEGEN_VER"

generate() {
	local k="$1" f="${1}_gen.go"
	shift

	trap "rm -f $f~" EXIT

	set -- -package "$GOPACKAGE" -generate "$k" $@

	echo "+ go run \"$OAPI_CODEGEN_URL\" $@ -o $f $OAPI_FILE" >&2
	go run "$OAPI_CODEGEN_URL" $@ -o "$f~" "$OAPI_FILE"

	if ! diff -u "$f" "$f~"; then
		mv "$f~" "$f"
	else
		rm "$f~"
	fi
}

generate types -alias-types
generate client
