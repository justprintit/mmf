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

type API struct {
	URL    string
	Method string `json:"httpMethod"`
}
