package util

import (
	"github.com/sethvargo/go-password/password"
)

func RandomString(n int) (string, error) {
	return password.Generate(n, n/4, n/8, false, true)
}
