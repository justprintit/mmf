package types

import (
	"fmt"
	"strconv"

	"github.com/justprintit/mmf/api/client/json"
)

type Id struct {
	s string
	n int
}

func NewId(v interface{}) (Id, error) {
	var w Id

	if s, ok := v.(string); ok {
		if w.SetFromString(s) {
			return w, nil
		}
	} else if n, ok := v.(int); ok {
		if w.SetInt(n) {
			return w, nil
		}
	}

	return w, ErrInvalidValue
}

func (a Id) Lt(b Id) bool {
	if len(a.s) > 0 {
		if len(b.s) > 0 {
			// string vs string
			return a.s < b.s
		} else {
			// string > number
			return false
		}
	} else if len(b.s) > 0 {
		// number < string
		return true
	} else {
		// number vs number
		return a.n < b.n
	}
}

func (w *Id) Ok() bool {
	if w != nil {
		if len(w.s) > 0 || w.n > 0 {
			return true
		}
	}
	return false
}

func (w *Id) Value() (interface{}, bool) {
	if w == nil {
		return nil, false
	} else if len(w.s) > 0 {
		return w.s, true
	} else {
		return w.n, w.n > 0
	}
}

func (w *Id) String() string {
	if w != nil {
		if len(w.s) > 0 {
			return w.s
		} else if w.n > 0 {
			return fmt.Sprintf("%v", w.n)
		}
	}
	return "INVALID"
}

func (w *Id) Int() (int, bool) {
	if w == nil || len(w.s) > 0 {
		return 0, false
	} else {
		return w.n, w.n > 0
	}
}

func (w *Id) SetInt(n int) bool {
	if w != nil && n > 0 {
		w.n = n
		w.s = ""
		return true
	}
	return false
}

func (w *Id) SetString(s string) bool {
	if w != nil && len(s) > 0 {
		w.n = 0
		w.s = s
		return true
	}
	return false
}

func (w *Id) SetFromString(s string) bool {
	if w == nil || len(s) == 0 {
		return false
	} else if n, err := strconv.Atoi(s); err == nil {
		return w.SetInt(n)
	} else {
		w.n = 0
		w.s = s
		return true
	}
}

func (w *Id) UnmarshalJSON(data []byte) error {
	if len(data) > 0 {
		if data[0] == '"' {
			var s string

			if err := json.Unmarshal(data, &s); err != nil {
				return err
			} else if w.SetFromString(s) {
				return nil
			}
		} else {
			var n int

			if err := json.Unmarshal(data, &n); err != nil {
				return err
			} else if w.SetInt(n) {
				return nil
			}
		}
	}

	return ErrInvalidValue
}

func (w *Id) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	} else if w.SetFromString(s) {
		return nil
	} else {
		return ErrInvalidValue
	}
}

func (w Id) MarshalJSON() ([]byte, error) {
	if v, ok := w.Value(); ok {
		return json.Marshal(v)
	}

	return []byte{}, ErrInvalidValue
}

func (w Id) MarshalYAML() (interface{}, error) {
	if v, ok := w.Value(); ok {
		return v, nil
	} else {
		return nil, ErrInvalidValue
	}
}
