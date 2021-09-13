package json

import (
	json "github.com/json-iterator/go"
)

type Groups struct {
	Count int     `json:"total_count"`
	Group []Group `json:"items"`
}

type Group struct {
	Id           GroupId        `json:"id"`
	Name         string         `json:"name"`
	API          map[string]API `json:"apis,omitempty"`
	TotalObjects int            `json:"total_count_objects,omitempty"`
	Children     []Group        `json:"childrens,omitempty"`
	Objects
}

type GroupId struct {
	id int
	s  string
}

func (w *GroupId) Int() (int, bool) {
	if len(w.s) > 0 {
		return 0, false
	} else {
		return w.id, true
	}
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
