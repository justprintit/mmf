package json

//go:generate ./api.sh

import (
	"github.com/go-resty/resty/v2"

	"github.com/justprintit/mmf/api/client"
)

// Code generated by ./api.sh DO NOT EDIT

// UserSharedGroupsResult() pulls a *Objects out of a resty.Response
func UserSharedGroupsResult(resp *resty.Response) *Objects {
	if out := resp.Result(); out != nil {
		return out.(*Objects)
	}
	return nil
}

// UserSharedLibraryResult() pulls a *UserSharedLibrary out of a resty.Response
func UserSharedLibraryResult(resp *resty.Response) *UserSharedLibrary {
	if out := resp.Result(); out != nil {
		return out.(*UserSharedLibrary)
	}
	return nil
}

// Pledges
var PledgesLibraryRequest = client.RequestOptions{
	Accept:  "application/json",
	Referer: "/library/?v=campaigns",
	Path:    "/data-library/campaigns",
	Method:  "GET",
	Result:  Objects{},
}

// PledgesLibraryResult() pulls a *Objects out of a resty.Response
func PledgesLibraryResult(resp *resty.Response) *Objects {
	if out := resp.Result(); out != nil {
		return out.(*Objects)
	}
	return nil
}

// Purchases
var PurchasesLibraryRequest = client.RequestOptions{
	Accept:  "application/json",
	Referer: "/library/?v=purchases",
	Path:    "/data-library/purchases",
	Method:  "GET",
	Result:  Objects{},
}

// PurchasesLibraryResult() pulls a *Objects out of a resty.Response
func PurchasesLibraryResult(resp *resty.Response) *Objects {
	if out := resp.Result(); out != nil {
		return out.(*Objects)
	}
	return nil
}

// Shared
var SharedLibraryRequest = client.RequestOptions{
	Accept:  "application/json",
	Referer: "/library/?v=shared",
	Path:    "/data-library/shared",
	Method:  "GET",
	Result:  Users{},
}

// SharedLibraryResult() pulls a *Users out of a resty.Response
func SharedLibraryResult(resp *resty.Response) *Users {
	if out := resp.Result(); out != nil {
		return out.(*Users)
	}
	return nil
}

// Tribes
var TribesLibraryRequest = client.RequestOptions{
	Accept:  "application/json",
	Referer: "/library/?v=my-tribes",
	Path:    "/data-library/tribes",
	Method:  "GET",
	Result:  Tribes{},
}

// TribesLibraryResult() pulls a *Tribes out of a resty.Response
func TribesLibraryResult(resp *resty.Response) *Tribes {
	if out := resp.Result(); out != nil {
		return out.(*Tribes)
	}
	return nil
}
