package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/juju/persistent-cookiejar"
	"github.com/motemen/go-loghttp"
	"golang.org/x/net/publicsuffix"

	"github.com/justprintit/mmf/api/transport"
	"github.com/justprintit/mmf/web/server"
)

const (
	CallbackPath = "/oauth2/callback"
)

type App struct {
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

func (m *App) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	http.NotFound(rw, req)
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

	m.server = srv
	m.worker = mmf
	return m, nil
}
