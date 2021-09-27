#!/bin/sh

set -eu

F="${0%.sh}_sh.go"
trap "rm -f '$F~'" EXIT
exec > "$F~"

LIBRARIES=
LIBRARIES="${LIBRARIES:+$LIBRARIES }Shared:Users"
LIBRARIES="${LIBRARIES:+$LIBRARIES }Purchases:Objects"
LIBRARIES="${LIBRARIES:+$LIBRARIES }Pledges:Objects"
LIBRARIES="${LIBRARIES:+$LIBRARIES }Tribes"

# Prelude
#
cat <<EOT
package library

//go:generate $0

import (
	"context"
	"os"

	"github.com/justprintit/mmf/api/library/json"
)

// Code generated by $0 DO NOT EDIT
EOT

# GetFooLibrary()
#
for x in $LIBRARIES; do
	n="${x%:*}"
	t="${x#*:}"

cat <<EOT

// get all JSON data of $n library
func (c *Client) Get${n}Library(ctx context.Context) (*json.$t, error) {
	d, err := c.Get${n}LibraryPage(ctx, 0)
	if err == nil {
		// pagination
		p := json.${t}Pages(d)

		// download further pages if needed
		page, _, ok := p.Next(1)
		for ok {
			var d2 *json.$t

			d2, err = c.Get${n}LibraryPage(ctx, page)
			if err != nil {
				break
			}

			d.Items = append(d.Items, d2.Items...)
		}
	}
	return d, err
}

// get requested page of JSON data of $n library
func (c *Client) Get${n}LibraryPage(ctx context.Context, page int) (*json.$t, error) {
	resp, err := c.GetPage(ctx, json.${n}LibraryRequest, page)
	if err != nil {
		os.Stdout.Write(resp.Body())
		return nil, err
	}

	out := json.${n}LibraryResult(resp)
	return out, nil
}
EOT
done

if [ -z "$LIBRARIES" ]; then
	cat <<EOT

func init() {
	// avoid "imported and not used" errors
	_ = context.Background
	_ = os.Open
	_ = json.Write
}
EOT
fi

if ! diff -u "$F" "$F~" >&2; then
	mv "$F~" "$F"
fi
