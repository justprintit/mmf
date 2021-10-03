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

SUBLIBRARIES=
SUBLIBRARIES="${SUBLIBRARIES:+$SUBLIBRARIES }UserSharedLibrary"
SUBLIBRARIES="${SUBLIBRARIES:+$SUBLIBRARIES }UserSharedGroup:"
SUBLIBRARIES="${SUBLIBRARIES:+$SUBLIBRARIES }TribeSharedGroup:"

# Prelude
#
cat <<EOT
package library

//go:generate $0

import (
	"context"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/justprintit/mmf/api/library/json"
	"github.com/justprintit/mmf/api/library/types"
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
cat <<EOT
	c.SchedulePageRequest(json.${n}LibraryRequest, 0, refresh${n}LibraryCallback)
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
	[ -n "$t" ] || t=Objects

cat <<EOT

// pull first page of $n library
func (c *Client) Refresh${n}Library(ctx context.Context) error {
	c.SchedulePageRequest(json.${n}LibraryRequest, 0, refresh${n}LibraryCallback)
	return nil
}

// handle first page of $n library
func refresh${n}LibraryCallback(c *Client, ctx context.Context, resp *resty.Response) error {
	d := json.${n}LibraryResult(resp)

	// pagination
	p := json.${t}Pages(d)

	// schedule further pages if needed
	page, offset, ok := p.Next(1)
	for ok {
		off := offset
		c.SchedulePageRequest(json.${n}LibraryRequest, page, func(c *Client, ctx context.Context, resp *resty.Response) error {
			d := json.${n}LibraryResult(resp)
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

# responses that require parsing the URL.Path
#
for x in $SUBLIBRARIES; do
	n="${x%:*}"
	t="${x#*:}"
	[ -n "$t" ] || t=Objects

	cat <<EOT

// handle first page of $n
func refresh${n}Callback(c *Client, ctx context.Context, resp *resty.Response) error {
	if d := json.${n}Result(resp); d != nil {
		req := resp.RawResponse.Request

		// process the first page
		if err := c.refresh${n}FromRequest(ctx, req, d); err != nil {
			return err
		}

		// and schedule further pages if needed
		p := json.${t}Pages(d)
		if page, _, ok := p.Next(1); ok {
			opt := json.New${n}FromRequest(req)

			for ok {
				c.SchedulePageRequest(opt, page, func(c *Client, ctx context.Context, resp *resty.Response) error {
					d := json.${n}Result(resp)
					req := resp.RawResponse.Request

					return c.refresh${n}FromRequest(ctx, req, d)
				})

				// next page
				page, _, ok = p.Next(page)
			}
		}
	}
	return nil
}
EOT
done

generate_schedule_l2() {
	local n="$1" t="$2" R="${3:-}" Rt="${4:-$t}"
	local va= vb=

	va="${t% *}"
	t="${t#* }"
	if [ -n "$R" ]; then
	       vb="${R% *}"
	       R="${R#* }"
	fi

	cat <<EOT

func (c *Client) schedule${n}(ctx context.Context, $va *$t) error {
	select {
	case <-ctx.Done():
		// cancelled
		return ctx.Err()
	default:
		t := time.Now()
EOT
		if [ -n "$R" ]; then
			# for each in range
			cat <<EOT
		for _, $vb := range $va.$R {
			if t.After($vb.Next${Rt}Update) {
				$vb.Next${Rt}Update = t.Add(Next${Rt}Update)

				req := json.New${n}Request($vb)
				c.Download(req, refresh${n}Callback)
			}
		}
EOT
		else
			# direct
			cat <<EOT
		if t.After($va.Next${n}Update) {
			$va.Next${n}Update = t.Add(Next${n}Update)

			req := json.New${n}Request($va)
			c.Download(req, refresh${n}Callback)
		}
EOT
		fi
	cat <<EOT
		return nil
	}
}
EOT
}

generate_schedule_l2 UserSharedLibrary "u types.User"
generate_schedule_l2 UserSharedGroup "u types.User" "g GroupsAll()" "GroupObjects"
generate_schedule_l2 TribeSharedGroup "u types.Tribe" "g GroupsAll()" "GroupObjects"

if [ -z "$LIBRARIES" ]; then
	cat <<EOT

func init() {
	// avoid "imported and not used" errors
	_ = resty.New
	_ = json.Write
}
EOT
fi

if ! diff -u "$F" "$F~" >&2; then
	mv "$F~" "$F"
fi
