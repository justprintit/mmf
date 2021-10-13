package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/juju/persistent-cookiejar"
	"github.com/motemen/go-loghttp"
	"golang.org/x/net/publicsuffix"

	"go.sancus.dev/web/errors"
	"go.sancus.dev/web/router"

	"github.com/justprintit/mmf/api/transport"
	"github.com/justprintit/mmf/web/server"
)

const (
	RedirectPath = "/oauth2"
	CallbackPath = "/oauth2/callback"
)

type App struct {
	http.Handler

	config     Config
	configFile string

	server *server.Server
	worker *transport.Client
}

func (m *App) LogRequest(req *http.Request) {
	loghttp.DefaultLogRequest(req)
}

func (m *App) LogResponse(resp *http.Response) {
	loghttp.DefaultLogResponse(resp)
}

func (m *App) ErrorHandler(rw http.ResponseWriter, req *http.Request, err error) {
	errors.HandleError(rw, req, err)
}

func (m *App) URL() string {
	return fmt.Sprintf("http://%s/", m.server.Addr)
}

func (m *App) Run() {
	// launch server
	go m.server.Serve()

	// and run the worker
	m.worker.Run()
}

func (m *App) Wait() {
	// wait until the worker finishes
	m.worker.Wait()

	// and shutdown the web server
	m.server.Shutdown(context.Background())
}

func (m *App) Reload() error {
	return nil
}

func (m *App) Abort() {
	m.worker.Abort()
}

func NewApp(cfg Config, cfgFile string) (*App, error) {

	m := &App{
		config:     cfg,
		configFile: cfgFile,
	}

	// cookiejar
	jar, err := cookiejar.New(&cookiejar.Options{
		Filename:         cfg.Cookies,
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, err
	}

	// server
	srv, err := cfg.Server.NewServer(m)
	if err != nil {
		return nil, err
	}

	// router
	r := router.NewRouter(m.ErrorHandler)

	// worker
	rt := &loghttp.Transport{
		LogRequest:  m.LogRequest,
		LogResponse: m.LogResponse,
	}

	mmf, err := transport.NewClientWithOptions(
		transport.WithTransport(rt),
		transport.WithCookieJar(jar),
		transport.WithUser(cfg.Auth.User),
		transport.WithOauth2(cfg.Auth.Client, srv.URL().String(), CallbackPath),
	)

	if err != nil {
		defer srv.Close()
		return nil, err
	}

	// oauth2
	r.TryHandleFunc(RedirectPath, mmf.RedirectHandler)
	r.TryHandleFunc(CallbackPath, mmf.CallbackHandler)

	m.worker = mmf // mmf client
	m.server = srv // http server
	m.Handler = r  // router
	return m, nil
}
