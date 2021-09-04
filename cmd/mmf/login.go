package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/kr/pretty"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"

	"go.sancus.dev/web/errors"
	"golang.org/x/oauth2"

	"github.com/justprintit/mmf/api/auth"
)

const (
	CodeExchangeTimeout = time.Second
	RandomStateLength   = 16
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "logs into MMF",
	RunE: func(cmd *cobra.Command, args []string) error {

		s, err := NewServer()
		if err != nil {
			return errors.Wrap(err, "NewServer")
		}

		h := LoginManager{
			base: s.BaseURL(),
		}

		h.Setup(s, cfg.Auth)
		go s.Serve()

		u := h.LoginURL()
		err = browser.OpenURL(u)
		if err != nil {
			return errors.Wrap(err, "OpenURL: %q", u)
		}

		<-time.After(10 * time.Minute)
		return nil
	},
}

type LoginManager struct {
	config *oauth2.Config
	base   url.URL
}

func (h *LoginManager) Setup(mux Router, cfg auth.Config) {
	mux.HandleFunc("/", h.Index)
	mux.HandleFunc("/callback", h.Callback)

	oc, err := cfg.NewOauth2(h.CallbackURL())
	if err != nil {
		log.Fatal(err)
	}
	h.config = oc
}

func (h *LoginManager) LoginURL() string {
	return h.base.String()
}

func (h *LoginManager) CallbackURL() string {
	u := h.base
	u.Path = "/callback"
	return u.String()
}

func (h *LoginManager) Index(rw http.ResponseWriter, req *http.Request) {
	state, err := auth.RandomState(RandomStateLength)
	if err != nil {
		log.Fatal(err)
	}

	u := h.config.AuthCodeURL(state)
	// TODO: store state in cookie

	// redirect to AuthCodeURL
	errors.NewSeeOther(u).ServeHTTP(rw, req)
}

func (h *LoginManager) Callback(rw http.ResponseWriter, req *http.Request) {
	var check errors.BadRequestError

	q := req.URL.Query()
	state := q.Get("state")

	if len(state) == RandomStateLength {
		// TODO: validate state against cookie

		code := q.Get("code")
		ctx, cancel := context.WithTimeout(req.Context(), CodeExchangeTimeout)
		defer cancel()

		token, err := h.config.Exchange(ctx, code)
		if err == nil {
			fmt.Fprintf(rw, "token:%# v\n", pretty.Formatter(token))
			return
		}
		check.AppendWrapped(err, "%s.%s", "oauth2", "Exchange")
	} else {
		check.InvalidArgument("state")
	}

	check.ServeHTTP(rw, req)
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
