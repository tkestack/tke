package proxy

import (
	"bytes"
	gojson "encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"tkestack.io/tke/pkg/mesh/util/json"
	"tkestack.io/tke/pkg/util/log"
)

func New(opts ...Opt) *Proxy {
	p := &Proxy{
		Request: &Request{},
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

type Proxy struct {
	// TargetUrl *url.URL

	// http request header to replace
	// Headers map[string][]string

	// http request json body to rewrite patches
	// Patches []json.Patch

	Request *Request

	Director       func(*http.Request)
	ModifyResponse func(*http.Response) error
	ErrorHandler   func(http.ResponseWriter, *http.Request, error)

	// proxy *httputil.ReverseProxy
}

type Request struct {
	Method      string
	URL         *url.URL
	Header      http.Header
	BodyPatches []json.Patch
	Body        []byte
}

func (p *Proxy) defaultDirector(req *http.Request) {
	if p.Request.URL == nil {
		p.Request.URL = req.URL
		return
	}
	if p.Request.Method == "" {
		p.Request.Method = req.Method
	}

	newReq := &http.Request{
		Method:     p.Request.Method,
		Host:       p.Request.URL.Host,
		URL:        p.Request.URL,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     req.Header,
		Body:       req.Body,
	}

	// replace with new request
	defer func() {
		*req = *newReq
	}()

	if p.Request.Header != nil {
		for key, value := range p.Request.Header {
			val := ""
			if len(value) == 1 {
				val = value[0]
			}
			newReq.Header.Set(key, val)
		}
	}

	// newReq.Header.Add("Content-Type", "application/json")
	if _, ok := newReq.Header["User-Agent"]; !ok {
		// explicitly disable User-Agent so it's not set to default value
		newReq.Header.Set("User-Agent", "")
	}

	if len(p.Request.Body) == 0 && req.Body != nil {
		err := p.Request.copyBody(req.Body)
		eMsg := "fleet-hub.proxy: read request body error, Cause: %s"
		if err != nil {
			log.Errorf(eMsg, err.Error())
			newReq.URL = nil
			return
		}
	}

	body, err := p.Request.modifyBody()
	if err != nil {
		eMsg := "fleet-hub.proxy: modify request body error, Cause: %s"
		log.Errorf(eMsg, err.Error())
		newReq.URL = nil
		return
	}
	if len(body) != 0 {
		newReq.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}
}

func (p *Proxy) Proxy(w http.ResponseWriter, req *http.Request) {
	director := p.defaultDirector
	if p.Director != nil {
		director = p.Director
	}

	proxy := &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: p.ModifyResponse,
		ErrorHandler:   p.ErrorHandler,
	}

	proxy.ServeHTTP(w, req)
}

func (r *Request) copyBody(body io.ReadCloser) error {
	var bodyBuf bytes.Buffer
	_, err := io.Copy(&bodyBuf, body)
	if err != nil {
		body.Close()
		return err
	}
	r.Body = bodyBuf.Bytes()
	return nil
}

func (r *Request) modifyBody() ([]byte, error) {
	if len(r.Body) != 0 && len(r.BodyPatches) != 0 {
		patchJSON, err := json.Marshal(r.BodyPatches)
		eMsg := "fleet-hub.proxy: modify request body error, Cause: %s"
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf(eMsg, err.Error()))
		}

		patch, err := jsonpatch.DecodePatch(patchJSON)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf(eMsg, err.Error()))
		}

		body, err := patch.Apply(r.Body)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf(eMsg, err.Error()))
		}

		return body, nil
	}

	return r.Body, nil
}

func NewDumpResponseWriter() *DumpResponseWriter {
	return &DumpResponseWriter{
		header: make(http.Header),
		Status: 0,
	}
}

type DumpResponseWriter struct {
	header http.Header

	Status int
	Body   bytes.Buffer
}

func (w *DumpResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(map[string][]string)
	}
	return w.header
}

func (w *DumpResponseWriter) Write(b []byte) (int, error) {
	return w.Body.Write(b)
}

func (w *DumpResponseWriter) WriteHeader(statusCode int) {
	w.Status = statusCode
}

func (w *DumpResponseWriter) Unmarshal(model BaseAPIResponse) (err error) {
	if err = json.Unmarshal(w.Body.Bytes(), model); err != nil {
		return errors.Wrap(err, fmt.Sprintf("fleet-hub.proxy: API json decode error, Cause: %s", err))
	}
	return
}

func (w *DumpResponseWriter) OnData(fn func(data jsoniter.Any)) (int, error) {
	bs := w.Body.Bytes()

	if w.Status != http.StatusOK {
		me := make([]string, 0)
		msg := json.Get(bs, "message")
		if msg.LastError() != nil {
			me = append(me, msg.LastError().Error())
		} else {
			me = append(me, msg.ToString())
		}
		e := json.Get(bs, "error")
		if e.LastError() != nil {
			me = append(me, e.LastError().Error())
		} else {
			me = append(me, e.ToString())
		}
		return w.Status, fmt.Errorf("%s", strings.Join(me, ", "))
	}

	data := json.Get(bs, "data")
	if data.LastError() != nil {
		log.Warnf(data.LastError().Error())
	} else {
		fn(data)
	}

	return w.Status, nil
}

func NewLoopPageProxy(f func(gojson.RawMessage) ([]interface{}, error), proxyOpts ...Opt) *LoopPageProxy {
	return &LoopPageProxy{
		proxy: New(proxyOpts...),
		ws:    make([]*DumpResponseWriter, 0),
		apiResponse: &APIResponse{
			ResultCode: 0,
			Msg:        "",
		},
		pageItems:     make([]interface{}, 0),
		unmarshalPage: f,
	}
}

type LoopPageProxy struct {
	proxy *Proxy

	ws []*DumpResponseWriter

	apiResponse *APIResponse

	pageItems []interface{}

	unmarshalPage func(gojson.RawMessage) ([]interface{}, error)
}

func (l *LoopPageProxy) LastWriter() DumpResponseWriter {
	if len(l.ws) == 0 {
		return DumpResponseWriter{}
	}
	return *l.ws[len(l.ws)-1]
}

func (l *LoopPageProxy) ResponseStatus() (status, resultCode int, msg string) {
	return l.LastWriter().Status, l.apiResponse.ResultCode, l.apiResponse.Msg
}

func (l *LoopPageProxy) Loop(req *http.Request) PageItems {

	var (
		query      = req.URL.Query()
		page       = 0
		count      = 1000
		resultCode int
		msg        string
		err        error
	)

	for ; ; page++ {
		query.Set("offset", strconv.Itoa(page*count))
		query.Set("limit", strconv.Itoa(count))
		req.URL.RawQuery = query.Encode()

		writer := NewDumpResponseWriter()
		l.proxy.Proxy(writer, req)
		l.ws = append(l.ws, writer)

		response := &APIPageResponse{}
		err = writer.Unmarshal(response)
		if err != nil {
			log.Errorf("fleet-hub.proxy: json decoding data error, Cause: %s", err.Error())
			break
		}

		resultCode = response.ResultCode
		msg = response.Msg

		if writer.Status != http.StatusOK {
			break
		}

		if l.unmarshalPage == nil {
			break
		}

		raw := response.RawData()
		if raw == nil {
			break
		}

		cur, err := l.unmarshalPage(*raw)
		if err != nil {
			log.Errorf("fleet-hub.proxy: json decoding page items error, Cause: %s", err.Error())
			break
		}
		if cur != nil {
			l.pageItems = append(l.pageItems, cur...)
		}

		if cur == nil || len(cur) < count {
			break
		}
	}

	l.apiResponse.ResultCode = resultCode
	l.apiResponse.Msg = msg
	return l.pageItems
}
