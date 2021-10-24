#!/bin/sh

set -eu
F="${0%.sh}_sh.go"
trap "rm -f '$F~'" EXIT
exec > "$F~"

cat <<EOT
package ${GOPACKAGE:-undefined}

//go:generate $0${1:+ $*}

import (
	"github.com/justprintit/mmf/api/openapi"
	"github.com/justprintit/mmf/util"
)

// prevent unused import errors
var _ = util.Pages
var _ = openapi.NewClient
EOT

generate() {
	local n="$1" x= k= t= cond=
	shift

# type
	cat <<EOT

// ${n}RequestParams are the parameters for $n requests
type ${n}RequestParams struct {
EOT
	for x; do
		k="${x%:*}"
		t="${x#*:}"
		cat <<EOT
	$k $t
EOT
	done
	cat <<EOT
}
EOT

# accessors
	for x; do
		k="${x%:*}"
		t="${x#*:}"

		case "$t" in
		string)
			#cond="rp.$k != \"\""
			cond="v != \"\""
			;;
		int)
			#cond="rp.$k > 0"
			cond="v > 0"
			;;
		*)
			cond=
			;;
		esac

		cat <<EOT

func (rp ${n}RequestParams) As${k}() openapi.${k} {
	return openapi.${k}(rp.$k)
}

func (rp ${n}RequestParams) As${k}Pointer() *openapi.${k} {
	if v := rp.As${k}(); $cond {
		return &v
	} else {
		return nil
	}
}
EOT
		if [ "PerPage" = "$k" ]; then
			cat<<EOT

func (rp ${n}RequestParams) Pages(total int) int {
	if rp.$k > 0 {
		return util.Pages(total, rp.$k, rp.$k)
	} else {
		return 0
	}
}
EOT
		fi
	done
}

generate User Username:string Page:int PerPage:int

# format
gofmt -w -l -s "$F~"

if ! diff -u "$F" "$F~" >&2; then
	mv "$F~" "$F"
fi
