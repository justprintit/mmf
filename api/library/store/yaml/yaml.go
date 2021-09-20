package yaml

import (
	"io"

	"gopkg.in/yaml.v2"
)

type (
	Decoder = yaml.Decoder
	Encoder = yaml.Encoder
)

func NewDecoder(in io.Reader) *yaml.Decoder {
	adapter := yaml.NewDecoder(in)
	adapter.SetStrict(true)
	return adapter
}

func NewEncoder(out io.Writer) *yaml.Encoder {
	adapter := yaml.NewEncoder(out)
	return adapter
}

func WriteTo(v interface{}, out io.Writer) (int64, error) {
	if b, err := yaml.Marshal(v); err != nil {
		return 0, err
	} else {
		if p, ok := out.(io.Closer); ok {
			defer p.Close()
		}
		n, err := out.Write(b)
		return int64(n), err
	}
}
