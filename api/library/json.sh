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

	"github.com/go-resty/resty/v2"

	"github.com/justprintit/mmf/api/library/json"
)

// Code generated by $0 DO NOT EDIT
EOT

# RefreshLibraries()
#
cat <<EOT

// pulls all data from the server
func (c *Client) RefreshLibraries(ctx context.Context) error {
	// load persistent data
	if err := c.Reload(); err != nil {
		return err
	}

	// schedule downloads to refresh libraries
EOT

for x in $LIBRARIES; do
	n="${x%:*}"
	t="${x#*:}"

	case "$n" in
	Pledges)
		g="campaigns"
		;;
	*)
		g="$(echo "$n" | tr 'A-Z' 'a-z')"
		;;
	esac

cat <<EOT
	c.ScheduleLibraryPageRequest("$g", 0, &json.$t{}, refresh${n}LibraryCallback)
EOT
done

cat <<EOT

	return nil
}
EOT

# RefreshFooLibrary()
#
for x in $LIBRARIES; do
	n="${x%:*}"
	t="${x#*:}"

	# library's name
	case "$n" in
	Pledges)
		g="campaigns"
		;;
	*)
		g="$(echo "$n" | tr A-Z a-z)"
		;;
	esac

	# Count handler
	case "$n" in
	Tribes)
		# int
		pages=Pages
		;;
	*)
		# json.Number
		pages=PagesN
		;;
	esac

cat <<EOT

// get first page of JSON data of $n library
func (c *Client) Get${n}Library(ctx context.Context) (*json.$t, error) {
	return c.Get${n}LibraryPage(ctx, 0)
}

// get requested page of JSON data of $n library
func (c *Client) Get${n}LibraryPage(ctx context.Context, page int) (*json.$t, error) {
	out := &json.$t{}

	resp, err := c.GetLibraryPage(ctx, "$g", page, out)
	if err != nil {
		os.Stdout.Write(resp.Body())
	}

	return out, err
}

// pull first page of $n library
func (c *Client) Refresh${n}Library(ctx context.Context) error {
	c.ScheduleLibraryRequest("$g", &json.$t{}, refresh${n}LibraryCallback)
	return nil
}

// handle first page of $n library
func refresh${n}LibraryCallback(c *Client, ctx context.Context, resp *resty.Response) error {
	d := resp.Result().(*json.$t)

	// pagination
	p := c.$pages(len(d.Items), d.Count)

	// schedule further pages if needed
	page, offset, ok := p.Next(1)
	for ok {
		off := offset
		c.ScheduleLibraryPageRequest("$g", page, &json.$t{}, func(c *Client, ctx context.Context, resp *resty.Response) error {
			d := resp.Result().(*json.$t)
			return c.refresh${n}Library(ctx, off, d.Items...)
		})

		// next page
		page, offset, ok = p.Next(page)
	}

	// and process first page
	return c.refresh${n}Library(ctx, 0, d.Items...)
}
EOT
done

if [ -z "$LIBRARIES" ]; then
	cat <<EOT

func init() {
	// avoid "imported and not used" errors
	_ = os.Open
	_ = resty.New
	_ = json.Write
}
EOT
fi

if ! diff -u "$F" "$F~" >&2; then
	mv "$F~" "$F"
fi