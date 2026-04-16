package proxy

import (
	"crypto/tls"
	"net/http"
	"sync"

	"github.com/elazarl/goproxy"
)

type AdBlockerProxy struct {
	server    *goproxy.ProxyHttpServer
	Blocklist *Blocklist
	enabled   bool
	mu        sync.RWMutex
}

func New(ca *tls.Certificate) *AdBlockerProxy {
	proxyServer := goproxy.NewProxyHttpServer()
	SetCAForGoproxy(proxyServer, ca)

	abProxy := &AdBlockerProxy{
		server:    proxyServer,
		Blocklist: NewBlocklist(),
		enabled:   true,
	}

	proxyServer.OnRequest().DoFunc(abProxy.OnRequest)
	return abProxy
}

func (p *AdBlockerProxy) OnRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	p.mu.RLock()
	active := p.enabled
	p.mu.RUnlock()

	if active && p.Blocklist.Contains(req.URL.Hostname()) {
		return req, goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusForbidden, "Blocked by AdBlocker")
	}
	return req, nil
}

func (p *AdBlockerProxy) Start(addr string) error {
	return http.ListenAndServe(addr, p.server)
}

func (p *AdBlockerProxy) SetEnabled(enabled bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.enabled = enabled
}

func (p *AdBlockerProxy) Enabled() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.enabled
}

func (p *AdBlockerProxy) UpdateBlocklist(urls []string) error {
	return p.Blocklist.UpdateFromURLs(urls)
}
