package gokinde

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
)

const WellKnownJWKsPath = "/.well-known/jwks"
const OAuth2TokenPath = "/oauth2/token"
const OAuth2AuthPath = "/oauth2/auth"
const UserProfilePath = "/oauth2/user_profile"

type Cfg struct {
	ClientID     string
	ClientSecret string

	KindeDomain string

	// Transport: optional, if nil a default transport with sensible rate limiting is used.
	Transport http.RoundTripper

	// ErrorLog: optional, if nil stdout is used
	ErrorLog func(err error)
}

type Client struct {
	jwks       *keyfunc.JWKS
	cfg        Cfg
	httpClient *http.Client
}

// NewClient initializes and returns a new client.
// The background thread is terminated when the context is cancelled.
// If this function returns an error - the client will not be function and your server should either:
//   - A: consider this a fatal error and panic to crash -- see MustStartAsyncWorker
//   - B: retry, by calling this function again (not recommended due to complexity, benefit and risk)
func NewClient(ctx context.Context, cfg Cfg) (*Client, error) {
	// TODO create rate limited http transport
	var trans http.RoundTripper

	cfg.KindeDomain = strings.TrimRight(cfg.KindeDomain, "/")

	cl := &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Transport: trans,
			Timeout:   60 * time.Second,
		},
	}
	return cl, cl.startAsyncWorker(ctx)
}

// MustNewClient calls NewClient and panics if it errors.
func MustNewClient(ctx context.Context, cfg Cfg) *Client {
	cl, err := NewClient(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return cl
}

// startAsyncWorker is non-blocking.
// The background thread is terminated when the context is cancelled.
// If this function returns an error - the client will not be function and your server should either:
//   - A: consider this a fatal error and panic to crash -- see MustStartAsyncWorker
//   - B: retry, by calling this function again (not recommended due to complexity, benefit and risk)
func (cl *Client) startAsyncWorker(ctx context.Context) error {
	if cl.cfg.ErrorLog == nil {
		cl.cfg.ErrorLog = func(err error) {
			log.Println("ERR: gokinde: fetching JWKs:", err)
		}
	}

	options := keyfunc.Options{
		Ctx:                 ctx,
		RefreshErrorHandler: cl.cfg.ErrorLog,
		RefreshInterval:     time.Hour,
		RefreshRateLimit:    time.Second * 30,
		RefreshTimeout:      time.Second * 45,
		RefreshUnknownKID:   true,
	}

	_, err := url.Parse(cl.cfg.KindeDomain)
	if err != nil {
		return fmt.Errorf("invalid KindeDomain: %w", err)
	}

	jwks, err := keyfunc.Get(cl.cfg.KindeDomain+WellKnownJWKsPath, options)
	if err != nil {
		return fmt.Errorf("fetching JWKs: %w", err)
	}

	cl.jwks = jwks

	return nil
}
