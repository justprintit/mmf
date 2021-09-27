package client

import (
	"github.com/justprintit/mmf/api/client/json"
)

type Pagination struct {
	Size  int
	Total int
}

func (p *Pagination) Next(prev int) (int, int, bool) {
	offset := prev * p.Size
	if offset < p.Total {
		return prev + 1, offset, true
	}
	return 0, 0, false
}

func (c *Client) Pages(n, total int) *Pagination {
	return &Pagination{
		Size:  n,
		Total: total,
	}
}

func (c *Client) PagesN(n int, count json.Number) *Pagination {
	total := n

	if n64, err := count.Int64(); err == nil {
		total = int(n64)
	}

	return c.Pages(n, total)
}
