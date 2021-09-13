package json

import (
	json "github.com/json-iterator/go"
)

type Users struct {
	Count json.Number `json:"total_count"`
	Items []User      `json:"items"`
}

type User struct {
	Id       string         `json:"id"`
	Username string         `json:"username"`
	Name     string         `json:"name"`
	Avatar   string         `json:"avatar_url"`
	API      map[string]API `json:"apis,omitempty"`
	Groups   Groups         `json:"groups,omitempty"`
}
