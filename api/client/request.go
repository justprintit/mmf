package client

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

const (
	DirectoryMode fs.FileMode = 0755
)

type ResponseHandler func(c *Client, req *Request, resp *resty.Response) error

type RequestOptions struct {
	Accept          string
	Referer         string
	Path            string
	Query           url.Values
	Method          string
	Result          interface{}
	OutputDirectory string
	Handler         ResponseHandler
}

func (opt RequestOptions) Clone() RequestOptions {
	// query
	q := make(map[string][]string, len(opt.Query))
	for k, v := range opt.Query {
		w := make([]string, len(v))
		copy(w, v)
		q[k] = w
	}

	opt.Query = q
	return opt
}

func (opt *RequestOptions) Setf(k string, v string, args ...interface{}) {
	if len(args) > 0 {
		v = fmt.Sprintf(v, args...)
	}
	opt.Query.Set(k, v)
}

func (opt RequestOptions) New(c *Client, ctx context.Context) *Request {
	// Request
	req := &Request{
		Request: c.Client.R(),
		Client:  c,
		Handler: opt.Handler,
	}

	// TraceInfo
	if c.TraceEnabled {
		req.EnableTrace()
	}

	// Context
	if ctx != nil {
		req.SetContext(ctx)
	}

	// Accept
	if s := opt.Accept; len(s) > 0 {
		req.SetHeader("Accept", s)
	}

	// Referer
	referer := opt.Referer
	if len(referer) == 0 {
		referer = c.HostURL
	} else if strings.HasPrefix(referer, "https://") || strings.HasPrefix(referer, "http://") {
		// ready
	} else if referer[0] != '/' {
		referer = fmt.Sprintf("%s/%s", c.HostURL, referer)
	} else {
		referer = c.HostURL + referer
	}
	req.SetHeader("Referer", referer)

	// URL
	path := opt.Path
	if len(path) == 0 {
		path = "/"
	} else if path[0] != '/' {
		path = "/" + path
	}

	s := make([]string, 2, 4)
	s[0] = c.HostURL
	s[1] = path

	if len(opt.Query) > 0 {
		s = append(s, "?", opt.Query.Encode())
	}
	req.URL = strings.Join(s, "")

	// Method
	method := opt.Method
	if len(method) == 0 {
		method = resty.MethodGet
	}
	req.Method = method

	// Result
	if opt.Result != nil {
		req.SetResult(opt.Result)
	} else if dir := opt.OutputDirectory; len(dir) > 0 {
		if err := os.MkdirAll(dir, DirectoryMode); err != nil {
			log.Println(err)
		}

		if f, err := os.CreateTemp(dir, ""); err != nil {
			log.Println(err)
		} else {
			req.Output = f
			req.Request.SetOutput(f.Name())
		}
	}

	return req
}

type Request struct {
	*resty.Request

	Client  *Client
	Handler ResponseHandler
	Output  *os.File
}

func (req *Request) Execute() error {
	if f := req.Output; f != nil {
		defer os.Remove(f.Name()) // clean up
	}

	if resp, err := req.Request.Send(); err != nil {
		return err
	} else {
		return req.Handler(req.Client, req, resp)
	}
}

func (req *Request) Get() (*resty.Response, error) {
	if f := req.Output; f != nil {
		defer os.Remove(f.Name()) // clean up
	}

	return req.Request.Send()
}

func (c *Client) R(referer string, args ...interface{}) *resty.Request {
	if len(args) > 0 {
		referer = fmt.Sprintf(referer, args...)
	}

	return RequestOptions{
		Referer: referer,
	}.New(c, nil).Request
}

func (c *Client) J(referer string, args ...interface{}) *resty.Request {
	req := c.R(referer, args...)
	req.SetHeader("Accept", "application/json")
	return req
}
