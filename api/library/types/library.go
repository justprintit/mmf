package types

import (
	"sync"
)

type Library struct {
	mu sync.Mutex

	User map[string]User `json:",omitempty"`
}
