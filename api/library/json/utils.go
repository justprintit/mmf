package json

import (
	"io"

	"github.com/justprintit/mmf/api/client/json"
)

func Write(data interface{}, indent string, out io.Writer) error {
	return json.Write(data, indent, out)
}
