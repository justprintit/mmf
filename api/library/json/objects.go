package json

import (
	"log"

	"go.sancus.dev/core/errors"

	"github.com/justprintit/mmf/api/client/json"
	"github.com/justprintit/mmf/api/library/types"
)

type ObjectItems []Object

func (w *ObjectItems) UnmarshalJSON(data []byte) error {
	var s []Object
	var err error

	if data[0] == '[' {
		err = json.Unmarshal(data, &s)
	} else {
		var m map[string]Object
		err = json.Unmarshal(data, &m)

		for _, v := range m {
			s = append(s, v)
		}
	}

	*w = s
	return err
}

type Objects struct {
	Count json.Number `json:"total_count,omitempty"`
	Items ObjectItems `json:"items,omitempty"`
}

type Object struct {
	Id                 int            `json:"id,omitempty"`
	Name               string         `json:"name,omitempty"`
	Description        string         `json:"description,omitempty"`
	Type               string         `json:"type,omitempty"`
	ObjType            string         `json:"document_name_s,omitempty"`
	Private            bool           `json:"is_private,omitempty"`
	Visits             int            `json:"visits,omitempty"`
	URL                string         `json:"url,omitempty"`
	ShowURL            string         `json:"show_url,omitempty"`
	AbsoluteURL        string         `json:"absolute_url,omitempty"`
	Image              string         `json:"obj_img,omitempty"`
	Images             Images         `json:"images,omitempty"`
	Wide               bool           `json:"wide,omitempty"`
	Purchased          bool           `json:"is_purchased,omitempty"`
	Price              ObjectPrice    `json:"price,omitempty"`
	FileMode           int            `json:"file_mode,omitempty"`
	Permissions        int            `json:"permissions,omitempty"`
	DownloadURL        string         `json:"download_url,omitempty"`
	ArchiveDownloadURL string         `json:"archive_download_url,omitempty"`
	Archives           []Archive      `json:"archives,omitempty"`
	Files              Files          `json:"files,omitempty"`
	Pledges            Groups         `json:"pledges,omitempty"`
	User               string         `json:"username,omitempty"`
	UserName           string         `json:"user_name,omitempty"`
	UserURL            string         `json:"user_url,omitempty"`
	UserImage          string         `json:"user_img,omitempty"`
	UserCollections    []Collection   `json:"user_collections,omitempty"`
	UserCredits        *json.Number   `json:"user_credits,omitempty"`
	API                map[string]API `json:"apis,omitempty"`
}

type ObjectPrice struct {
	Currency string      `json:"currency"`
	Symbol   string      `json:"symbol"`
	Value    json.Number `json:"value"`
}

func (obj *Object) Apply(d *types.Library, u *types.User, g *types.Group) error {
	log.Printf("%s[%v]: %q", obj.ObjType, obj.Id, obj.Name)
	return nil
}

func (w *Objects) Apply(d *types.Library, u *types.User, g *types.Group) error {
	var check errors.ErrorStack

	if n := len(w.Items); n > 0 {
		if v, err := w.Count.Int64(); err == nil {
			if int64(n) != v {
				log.Printf("Objects: expected:%v != actual:%v", v, n)
			}
		}

		for _, obj := range w.Items {
			if err := obj.Apply(d, u, g); err != nil {
				check.AppendError(err)
			}
		}
	}

	if err := check.AsError(); err != nil {
		d.OnError(u, err)
		return err
	}
	return nil
}
