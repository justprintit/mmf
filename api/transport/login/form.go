package login

import (
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Input struct {
	Name  string
	Type  string
	Value string
}

type Form struct {
	Id     string
	Action string
	Method string

	Inputs []Input
}

func (form *Form) Values() (url.Values, error) {
	v := url.Values{}

	for _, d := range form.Inputs {
		v.Add(d.Name, d.Value)
	}

	return v, nil
}

func (form *Form) AddInput(i int, s *goquery.Selection) {
	d := Input{}
	if v, ok := s.Attr("name"); ok {
		d.Name = v
	}
	if v, ok := s.Attr("type"); ok {
		d.Type = v
	}
	if v, ok := s.Attr("value"); ok {
		d.Value = v
	}

	if d.Name != "" {
		form.Inputs = append(form.Inputs, d)
	}
}

func (form *Form) EachInput(fn func(int, *Input)) {
	for i := range form.Inputs {
		d := &form.Inputs[i]
		fn(i, d)
	}
}

func parseForm(s *goquery.Selection) *Form {
	f := &Form{}
	if s, ok := s.Attr("id"); ok {
		f.Id = s
	}
	if s, ok := s.Attr("action"); ok {
		f.Action = s
	}
	if s, ok := s.Attr("method"); ok {
		f.Method = strings.ToUpper(s)
	}
	return f
}

func ParseForm(i int, s *goquery.Selection) *Form {
	if f := parseForm(s); f != nil {
		s.Find("input").Each(func(j int, s *goquery.Selection) {
			f.AddInput(j, s)
		})
		return f
	}
	return nil
}

func (form *Form) NewRequest(referer string) (*http.Request, error) {
	var (
		u0, u1   *url.URL
		body     string
		bodyType string
		err      error
	)

	// fields
	if data, err := form.Values(); err != nil {
		return nil, err
	} else {
		body = data.Encode()
		bodyType = "application/x-www-form-urlencoded"
	}

	// Referer
	if u0, err = url.Parse(referer); err != nil {
		return nil, err
	}

	// Location
	if form.Action == "" {
		u1 = u0
	} else if u1, err = url.Parse(form.Action); err != nil {
		return nil, err
	}

	if u1.Host != "" {
		// ready
	} else if u1.Path == "" {
		q := u1.RawQuery

		*u1 = *u0
		// restore Query if any
		if q != "" {
			u1.RawQuery = q
		}
	} else if u1.Path[0] == '/' {
		q := u1.RawQuery
		p := u1.Path

		*u1 = *u0
		// restore Query and Path
		u1.Path = p
		u1.RawQuery = q
	} else {
		q := u1.RawQuery
		p := u1.Path

		*u1 = *u0
		// join Path and restore Query
		u1.Path = path.Join(u0.Path, p)
		u1.RawQuery = q
	}

	req, err := http.NewRequest(form.Method, u1.String(), strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", bodyType)
	req.Header.Add("Content-Length", strconv.Itoa(len(body)))
	req.Header.Add("Referer", u0.String())
	return req, nil
}
