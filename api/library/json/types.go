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
	Id       json.Number
	Name     string
	Objects  int            `json:"total_count_objects,omitempty"`
	Children []Group        `json:"childrens,omitempty"`
	API      map[string]API `json:"apis,omitempty"`
}

type API struct {
	URL    string
	Method string `json:"httpMethod"`
}
