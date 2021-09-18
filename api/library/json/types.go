package json

import (
	"time"
)

type Archive struct {
	Id          int    `json:"id"`
	Filename    string `json:"path"`
	Size        int    `json:"size"`
	DownloadURL string `json:"download_url"`
}

type Files struct {
	Count int    `json:"total_count,omitempty"`
	Items []File `json:"items,omitempty"`
}

type File struct {
	Id           int       `json:"id"`
	Filename     string    `json:"filename"`
	Description  string    `json:"description"`
	Status       int       `json:"status"`
	StatusName   string    `json:"status_name"`
	Size         int       `json:"size"`
	PatchURL     string    `json:"patch_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	DownloadURL  string    `json:"download_url"`
	CreatedAt    time.Time `json:"created_at"`
	ViewerURL    string    `json:"viewer_url"`
	Render360    []string  `json:"render360_urls"`
}

type Images struct {
	Count int     `json:"total_count,omitempty"`
	Items []Image `json:"items,omitempty"`
}

type Image struct {
	Id                 int       `json:"id"`
	UploadId           string    `json:"upload_id"`
	Primary            bool      `json:"is_primary,omitempty"`
	PrintImageSelected bool      `json:"is_print_image_selected,omitempty"`
	Original           ImageFile `json:"original,omitempty"`
	Tiny               ImageFile `json:"tiny,omitempty"`
	Thumbnail          ImageFile `json:"thumbnail,omitempty"`
	Standard           ImageFile `json:"standard,omitempty"`
	Large              ImageFile `json:"large,omitempty"`
}

type ImageFile struct {
	URL    string `json:"url"`
	Width  *int   `json:"width"`
	Height *int   `json:"height"`
}

type Collection struct{}

type API struct {
	URL    string `json:"url"`
	Method string `json:"httpMethod"`
}
