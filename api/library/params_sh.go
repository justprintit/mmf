package library

//go:generate ./params.sh

import (
	"github.com/justprintit/mmf/api/openapi"
)

// UserRequestParams are the parameters for User requests
type UserRequestParams struct {
	Username string
	Page     int
	PerPage  int
}

func (rp UserRequestParams) AsUsername() openapi.Username {
	return openapi.Username(rp.Username)
}

func (rp UserRequestParams) AsUsernamePointer() *openapi.Username {
	if v := rp.AsUsername(); v != "" {
		return &v
	} else {
		return nil
	}
}

func (rp UserRequestParams) AsPage() openapi.Page {
	return openapi.Page(rp.Page)
}

func (rp UserRequestParams) AsPagePointer() *openapi.Page {
	if v := rp.AsPage(); v > 0 {
		return &v
	} else {
		return nil
	}
}

func (rp UserRequestParams) AsPerPage() openapi.PerPage {
	return openapi.PerPage(rp.PerPage)
}

func (rp UserRequestParams) AsPerPagePointer() *openapi.PerPage {
	if v := rp.AsPerPage(); v > 0 {
		return &v
	} else {
		return nil
	}
}