package auth

import (
	"github.com/sethvargo/go-password/password"
)

func RandomState(n int) (string, error) {
	return password.Generate(n, n/4, n/8, false, true)
}
