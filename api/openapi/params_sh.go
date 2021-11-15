package openapi

//go:generate ./params.sh

import (
	"github.com/justprintit/mmf/util"
)

// prevent unused import errors
var _ = util.Pages

// UserRequestParams are the parameters for User requests
type UserRequestParams struct {
	Username string
	Page     int
	PerPage  int
}

func (rp UserRequestParams) AsUsername() Username {
	return Username(rp.Username)
}

func (rp UserRequestParams) AsUsernamePointer() *Username {
	if v := rp.AsUsername(); v != "" {
		return &v
	} else {
		return nil
	}
}

func (rp UserRequestParams) AsPage() Page {
	return Page(rp.Page)
}

func (rp UserRequestParams) AsPagePointer() *Page {
	if v := rp.AsPage(); v > 0 {
		return &v
	} else {
		return nil
	}
}

func (rp UserRequestParams) AsPerPage() PerPage {
	return PerPage(rp.PerPage)
}

func (rp UserRequestParams) AsPerPagePointer() *PerPage {
	if v := rp.AsPerPage(); v > 0 {
		return &v
	} else {
		return nil
	}
}

func (rp UserRequestParams) Pages(total int) int {
	if rp.PerPage > 0 {
		return util.Pages(total, rp.PerPage, rp.PerPage)
	} else {
		return 0
	}
}
