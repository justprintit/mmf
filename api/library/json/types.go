package json

import (
	json "github.com/json-iterator/go"
)

type Users struct {
	Count int    `json:"total_count"`
	User  []User `json:"items"`
}

type User struct {
	Id       string
	Username string
	Name     string
	Avatar   string         `json:"avatar_url"`
	API      map[string]API `json:"apis,omitempty"`
	Groups   Groups
}

type Groups struct {
	Count int     `json:"total_count"`
	Group []Group `json:"items"`
}

type Group struct {
	Id       GroupId
	Name     string
	Objects  int            `json:"total_count_objects,omitempty"`
	Count    json.Number    `json:"total_count,omitempty"`
	Children []Group        `json:"childrens,omitempty"`
	Items    GroupItems     `json:",omitempty"`
	API      map[string]API `json:"apis,omitempty"`
}

type GroupId struct {
	id int
	s  string
}

func (w *GroupId) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		w.id = 0
		return json.Unmarshal(data, &w.s)
	} else {
		w.s = ""
		return json.Unmarshal(data, &w.id)
	}
}

func (w *GroupId) MarshalJSON() ([]byte, error) {
	if len(w.s) > 0 {
		return json.Marshal(w.s)
	} else {
		return json.Marshal(w.id)
	}
}

type GroupItems struct {
	s []Object
	m map[string]Object
}

func (w *GroupItems) UnmarshalJSON(data []byte) error {
	if data[0] == '[' {
		w.m = nil
		return json.Unmarshal(data, &w.s)
	} else {
		w.s = nil
		return json.Unmarshal(data, &w.m)
	}
}

func (w *GroupItems) MarshalJSON() ([]byte, error) {
	if w.m == nil {
		return json.Marshal(w.s)
	} else {
		return json.Marshal(w.m)
	}
}

type Objects struct {
	Count json.Number `json:"total_count,omitempty"`
	Items []Object
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

type Archive struct {
	Id          int
	Filename    string `json:"path"`
	Size        int    ``
	DownloadURL string `json:"download_url"`
}

type Files struct {
	Count int    `json:"total_count,omitempty"`
	Items []File `json:",omitempty"`
}

type File struct{}

type Images struct {
	Count int     `json:"total_count,omitempty"`
	Items []Image `json:",omitempty"`
}

type Image struct {
	Id                 int
	UploadId           string    `json:"upload_id"`
	Primary            bool      `json:"is_primary,omitempty"`
	PrintImageSelected bool      `json:"is_print_image_selected,omitempty"`
	Original           ImageFile `json:",omitempty"`
	Tiny               ImageFile `json:",omitempty"`
	Thumbnail          ImageFile `json:",omitempty"`
	Standard           ImageFile `json:",omitempty"`
	Large              ImageFile `json:",omitempty"`
}

type ImageFile struct {
	URL    string
	Width  *json.Number
	Height *json.Number
}

type Collection struct{}

type API struct {
	URL    string
	Method string `json:"httpMethod"`
}
