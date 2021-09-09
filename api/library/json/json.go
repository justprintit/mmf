package json

import (
	"bytes"
	"io"

	json "github.com/json-iterator/go"
)

type (
	Decoder = json.Decoder
	Encoder = json.Encoder
)

func NewDecoderBytes(body []byte) *json.Decoder {
	return NewDecoder(bytes.NewReader(body))
}

func NewDecoder(in io.Reader) *json.Decoder {
	adapter := json.NewDecoder(in)
	adapter.DisallowUnknownFields() // strict
	return adapter
}

func NewEncoder(out io.Writer) *json.Encoder {
	return json.NewEncoder(out)
}

func Write(data interface{}, indent string, out io.Writer) error {
	adapter := NewEncoder(out)
	if len(indent) > 0 {
		adapter.SetIndent("", indent)
	}
	return adapter.Encode(data)
}
