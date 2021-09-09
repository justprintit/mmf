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
	Children []Group        `json:"childrens,omitempty"`
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

type Objects struct {
	Count json.Number `json:"total_count,omitempty"`
	Items []Object
}

type Object struct {
	Id              int
	Name            string
	Type            string       `json:",omitempty"`
	ObjType         string       `json:"document_name_s,omitempty"`
	Private         bool         `json:"is_private,omitempty"`
	Visits          int          `json:",omitempty"`
	ShowURL         string       `json:"show_url,omitempty"`
	AbsoluteURL     string       `json:"absolute_url,omitempty"`
	Image           string       `json:"obj_img,omitempty"`
	Wide            bool         `json:",omitempty"`
	Purchased       bool         `json:"is_purchased,omitempty"`
	Price           ObjectPrice  `json:",omitempty"`
	DownloadURL     string       `json:"download_url,omitempty"`
	Archives        []Archive    `json:",omitempty"`
	User            string       `json:"username,omitempty"`
	UserName        string       `json:"user_name,omitempty"`
	UserURL         string       `json:"user_url,omitempty"`
	UserImage       string       `json:"user_img,omitempty"`
	UserCollections []Collection `json:"user_collections,omitempty"`
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

type Collection struct{}

type API struct {
	URL    string
	Method string `json:"httpMethod"`
}
