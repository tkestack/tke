package proxy

import (
	"net/http"
	"net/url"

	"tkestack.io/tke/pkg/mesh/util/json"
)

type Opt func(*Proxy)

func WithPatch(pch ...json.Patch) Opt {
	return func(p *Proxy) {
		p.Request.BodyPatches = pch
	}
}

func WithHeader(h http.Header) Opt {
	return func(p *Proxy) {
		p.Request.Header = h
	}
}

func WithRequestMethod(method string) Opt {
	return func(p *Proxy) {
		p.Request.Method = method
	}
}

func WithModifyResponse(f func(*http.Response) error) Opt {
	return func(p *Proxy) {
		p.ModifyResponse = f
	}
}

func WithTargetURL(u *url.URL) Opt {
	return func(p *Proxy) {
		p.Request.URL = u
	}
}
