package json

import (
	json "github.com/json-iterator/go"
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
	Id                 int
	Name               string
	Description        string         `json:",omitempty"`
	Type               string         `json:",omitempty"`
	ObjType            string         `json:"document_name_s,omitempty"`
	Private            bool           `json:"is_private,omitempty"`
	Visits             int            `json:",omitempty"`
	URL                string         `json:",omitempty"`
	ShowURL            string         `json:"show_url,omitempty"`
	AbsoluteURL        string         `json:"absolute_url,omitempty"`
	Image              string         `json:"obj_img,omitempty"`
	Images             Images         `json:",omitempty"`
	Wide               bool           `json:",omitempty"`
	Purchased          bool           `json:"is_purchased,omitempty"`
	Price              ObjectPrice    `json:",omitempty"`
	FileMode           int            `json:"file_mode,omitempty"`
	Permissions        int            `json:"permissions,omitempty"`
	DownloadURL        string         `json:"download_url,omitempty"`
	ArchiveDownloadURL string         `json:"archive_download_url,omitempty"`
	Archives           []Archive      `json:",omitempty"`
	Files              Files          `json:",omitempty"`
	Pledges            Groups         `json:",omitempty"`
	User               string         `json:"username,omitempty"`
	UserName           string         `json:"user_name,omitempty"`
	UserURL            string         `json:"user_url,omitempty"`
	UserImage          string         `json:"user_img,omitempty"`
	UserCollections    []Collection   `json:"user_collections,omitempty"`
	UserCredits        *json.Number   `json:"user_credits,omitempty"`
	API                map[string]API `json:"apis,omitempty"`
}

type ObjectPrice struct {
	Currency string
	Symbol   string
	Value    json.Number
}
