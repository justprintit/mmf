package json

import (
	"bytes"
	"io"

	json "github.com/json-iterator/go"
)

type (
	Number  = json.Number
	Config  = json.Config
	Decoder = json.Decoder
	Encoder = json.Encoder
)

var ConfigDefault = json.Config{
	// strict
	DisallowUnknownFields:         true,
	CaseSensitive:                 true,
	ObjectFieldMustBeSimpleString: true,
	ValidateJsonRawMessage:        true,
}.Froze()

func Marshal(v interface{}) ([]byte, error) {
	return ConfigDefault.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return ConfigDefault.Unmarshal(data, v)
}

func NewDecoderBytes(body []byte) *json.Decoder {
	return NewDecoder(bytes.NewReader(body))
}

func NewDecoder(in io.Reader) *json.Decoder {
	return ConfigDefault.NewDecoder(in)
}

func NewEncoder(out io.Writer) *json.Encoder {
	return ConfigDefault.NewEncoder(out)
}

func Write(data interface{}, indent string, out io.Writer) error {
	adapter := NewEncoder(out)
	adapter.SetIndent("", indent)
	return adapter.Encode(data)
}

func WriteTo(out io.Writer, data interface{}) error {
	return Write(data, "  ", out)
}
