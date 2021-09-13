package json

type Archive struct {
	Id          int
	Filename    string `json:"path"`
	Size        int    ``
	DownloadURL string `json:"download_url"`
}

type Files struct {
	Count int    `json:"total_count,omitempty"`
	Items []File `json:",omitempty"`
}

type File struct{}

type Images struct {
	Count int     `json:"total_count,omitempty"`
	Items []Image `json:",omitempty"`
}

type Image struct {
	Id                 int
	UploadId           string    `json:"upload_id"`
	Primary            bool      `json:"is_primary,omitempty"`
	PrintImageSelected bool      `json:"is_print_image_selected,omitempty"`
	Original           ImageFile `json:",omitempty"`
	Tiny               ImageFile `json:",omitempty"`
	Thumbnail          ImageFile `json:",omitempty"`
	Standard           ImageFile `json:",omitempty"`
	Large              ImageFile `json:",omitempty"`
}

type ImageFile struct {
	URL    string
	Width  *int
	Height *int
}

type Collection struct{}

type API struct {
	URL    string `json:"url"`
	Method string `json:"httpMethod"`
}
