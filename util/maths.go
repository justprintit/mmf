package util

import (
	"log"
)

func Pages(total, count, per_page int) int {
	pages := (total + per_page - 1) / per_page
	log.Printf("pages: %v/%v (%v) -> %v", count, total, per_page, pages)
	return pages
}
