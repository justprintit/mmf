package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync"

	"github.com/juju/persistent-cookiejar"
	"github.com/motemen/go-loghttp"
	"github.com/tidwall/pretty"
	"golang.org/x/net/publicsuffix"

	"go.sancus.dev/web/errors"
	"go.sancus.dev/web/router"

	"github.com/justprintit/mmf/api/library"
	"github.com/justprintit/mmf/api/library/store/bolt"
	"github.com/justprintit/mmf/api/transport"
	"github.com/justprintit/mmf/types"
	"github.com/justprintit/mmf/web/server"
)

const (
	RedirectPath = "/oauth2"
	CallbackPath = "/oauth2/callback"

	DownloadThreads = 3
)

type App struct {
	http.Handler

	config        Config
	configFile    string
	dumpTransport io.Writer

	mu     sync.Mutex
	data   types.Store
	server *server.Server
	worker *library.Worker
}

func (m *App) LogRequest(req *http.Request) {
	loghttp.DefaultLogRequest(req)

	if out := m.dumpTransport; out == nil {
		//
	} else if b, err := httputil.DumpRequest(req, true); err == nil {
		out.Write(b)
		out.Write([]byte{'\n', '\n'})
	} else {
		log.Println(err)
	}
}

func (m *App) LogResponse(resp *http.Response) {
	loghttp.DefaultLogResponse(resp)

	if out := m.dumpTransport; out != nil {
		var json bool

		switch {
		case strings.Contains(resp.Header.Get("Content-Type"), "json"):
			json = true
		}

		if b, err := httputil.DumpResponse(resp, !json); err == nil {
			out.Write(b)
			if json {
				// read and unread
				body, _ := ioutil.ReadAll(resp.Body)
				resp.Body.Close()
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

				out.Write(pretty.Pretty(body))
			}
			out.Write([]byte{'\n', '\n'})
		} else {
			log.Println(err)
		}
	}
}

func (m *App) ErrorHandler(rw http.ResponseWriter, req *http.Request, err error) {
	errors.HandleError(rw, req, err)
}

func (m *App) URL() string {
	return fmt.Sprintf("http://%s/", m.server.Addr)
}

func (m *App) updateConfig(v *string, after string, s string, args ...interface{}) bool {
	before := *v
	if before != after {
		if len(args) > 0 {
			s = fmt.Sprintf(s, args...)
		}
		if len(before) > 0 {
			log.Printf("%c %s: %q", '-', s, before)
		}
		log.Printf("%c %s: %q", '+', s, after)
		*v = after
		return true
	}
	return false
}

func (m *App) save() error {
	var err error

	if m.configFile == "" {
		_, err = m.config.WriteTo(os.Stdout)
	} else {
		_, err = m.config.WriteFile(m.configFile)
	}
	return err
}

func (m *App) onNewCredentials(user, password string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	f1 := m.updateConfig(&m.config.Auth.User.Username, user, "%s.%s.%s", "auth", "user", "username")
	f2 := m.updateConfig(&m.config.Auth.User.Password, password, "%s.%s.%s", "auth", "user", "password")
	if f1 || f2 {
		return m.save()
	}
	return nil
}

func (m *App) onNewClient(key, secret string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	f1 := m.updateConfig(&m.config.Auth.Client.ClientID, key, "%s.%s.%s", "auth", "api", "client_key")
	f2 := m.updateConfig(&m.config.Auth.Client.ClientSecret, secret, "%s.%s.%s", "auth", "api", "client_secret")
	if f1 || f2 {
		return m.save()
	}
	return nil
}

func (m *App) onNewToken(access, renew string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	f1 := m.updateConfig(&m.config.Auth.Client.AccessToken, access, "%s.%s.%s", "auth", "api", "access_token")
	f2 := m.updateConfig(&m.config.Auth.Client.RefreshToken, renew, "%s.%s.%s", "auth", "api", "refresh_token")
	if f1 || f2 {
		return m.save()
	}
	return nil
}

func (m *App) Run() {
	// launch server
	go m.server.Serve()

	// schedule some work
	go m.worker.Refresh()

	// and run the worker
	m.worker.Run(nil, DownloadThreads)
}

func (m *App) Wait() {
	// wait until the worker finishes
	m.worker.Wait()

	// and shutdown the web server
	m.server.Shutdown(context.Background())
}

func (m *App) Reload() error {
	return m.worker.Refresh()
}

func (m *App) Abort() {
	m.worker.Cancel()
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

	// library
	data, err := bolt.New(cfg.Data, types.LibraryEvents{})
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
		transport.WithCallbacks(transport.ClientEvents{
			OnNewCredentials: m.onNewCredentials,
			OnNewClient:      m.onNewClient,
			OnNewToken:       m.onNewToken,
		}),
	)
	if err != nil {
		defer srv.Close()
		return nil, err
	}

	w := library.NewWorker(mmf, data)

	// oauth2
	r.TryHandleFunc(RedirectPath, mmf.RedirectHandler)
	r.TryHandleFunc(CallbackPath, mmf.CallbackHandler)

	m.data = data  // library
	m.worker = w   // mmf client
	m.server = srv // http server
	m.Handler = r  // router

	// and start preloading data
	go data.Load()

	return m, nil
}
