package types

import (
	"log"
)

func UpdateString(label string, v *string, next string) error {
	if prev := *v; prev != next {
		if len(prev) > 0 {
			log.Printf("- %s: %q", label, prev)
		}
		log.Printf("+ %s: %q", label, next)
		*v = next
	}
	return nil
}
