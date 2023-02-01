package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/docker/distribution/manifest/manifestlist"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/registry/auth/token"
	"github.com/docker/libtrust"
	"tkestack.io/tke/api/registry"
	"tkestack.io/tke/pkg/registry/distribution/auth"
	"tkestack.io/tke/pkg/util/log"
)

var ManifestAccepts = []string{
	manifestlist.MediaTypeManifestList,
	schema2.MediaTypeManifest,
	schema1.MediaTypeSignedManifest,
	schema1.MediaTypeManifest,
}

// Repository holds information of a repository entity
type Repository struct {
	Endpoint   *url.URL
	client     *http.Client
	privateKey libtrust.PrivateKey
}

// NewRepository returns an instance of Repository
func NewRepository(endpoint string, privateKey libtrust.PrivateKey) (*Repository, error) {
	u, err := ParseEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	repository := &Repository{
		Endpoint:   u,
		client:     client,
		privateKey: privateKey,
	}

	return repository, nil
}

// ParseEndpoint parses endpoint to a URL
func ParseEndpoint(endpoint string) (*url.URL, error) {
	endpoint = strings.Trim(endpoint, " ")
	endpoint = strings.TrimRight(endpoint, "/")
	if len(endpoint) == 0 {
		return nil, fmt.Errorf("empty URL")
	}
	i := strings.Index(endpoint, "://")
	if i >= 0 {
		scheme := endpoint[:i]
		if scheme != "http" && scheme != "https" {
			return nil, fmt.Errorf("invalid scheme: %s", scheme)
		}
	} else {
		endpoint = "http://" + endpoint
	}

	return url.ParseRequestURI(endpoint)
}

// DeleteTag ...
func (r *Repository) DeleteTag(repoName, tag, user, tenantID string) error {
	digest, exist, err := r.ManifestExist(tag, repoName, tag, user, tenantID)
	if err != nil {
		return err
	}

	if !exist {
		log.Warnf("repo: %s:%s manifests not found.", repoName, tag)
		return nil
	}

	if err := r.DeleteManifest(digest, repoName, tag, user, tenantID); err != nil {
		return err
	}
	return nil
}

// ListTag ...
func (r *Repository) ListTag(repoName, user, tenantID string) ([]string, error) {
	tags := []string{}
	req, err := http.NewRequest("GET", buildTagListURL(r.Endpoint.String(), repoName), nil)
	if err != nil {
		return tags, err
	}
	err = r.withAuthInfo(req, repoName, user, tenantID)
	if err != nil {
		return tags, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return tags, parseError(err)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tags, err
	}

	if resp.StatusCode == http.StatusOK {
		tagsResp := struct {
			Tags []string `json:"tags"`
		}{}

		if err := json.Unmarshal(b, &tagsResp); err != nil {
			return tags, err
		}
		sort.Strings(tags)
		tags = tagsResp.Tags

		return tags, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return tags, nil
	}

	return tags, &Error{
		Code:    resp.StatusCode,
		Message: string(b),
	}

}

// ManifestExist ...
func (r *Repository) ManifestExist(reference, repoName, tag, user, tenantID string) (digest string, exist bool, err error) {
	req, err := http.NewRequest("HEAD", buildManifestURL(r.Endpoint.String(), repoName, reference), nil)
	if err != nil {
		return
	}
	err = r.withAuthInfo(req, repoName, user, tenantID)
	if err != nil {
		return
	}

	for _, mediaType := range ManifestAccepts {
		req.Header.Add("Accept", mediaType)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		err = parseError(err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		exist = true
		digest = resp.Header.Get("Docker-Content-Digest")
		return
	}

	if resp.StatusCode == http.StatusNotFound {
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = &Error{
		Code:    resp.StatusCode,
		Message: string(b),
	}
	return
}

// DeleteManifest ...
func (r *Repository) DeleteManifest(digest, repoName, tag, user, tenantID string) error {
	req, err := http.NewRequest("DELETE", buildManifestURL(r.Endpoint.String(), repoName, digest), nil)
	if err != nil {
		return err
	}
	err = r.withAuthInfo(req, repoName, user, tenantID)
	if err != nil {
		return err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return parseError(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted {
		return nil
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return &Error{
		Code:    resp.StatusCode,
		Message: string(b),
	}
}

func (r *Repository) withAuthInfo(req *http.Request, repoName, user, tenantID string) error {
	access := []*token.ResourceActions{
		{
			Type:    "repository",
			Actions: []string{"*", "pull"},
			// to make token be available, should rename repo name with tenantID
			Name: fmt.Sprintf("%s-%s", tenantID, repoName),
		},
	}
	token, err := auth.MakeToken(user, access, 24, r.privateKey)
	if err != nil {
		return err
	}
	log.Infof("token: %s", token.GetToken())
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.GetToken()))
	// set registry client UA to avoid error cased by reporting event of the deleted repo
	req.Header.Set("User-Agent", registry.RegistryClientUserAgent)
	return nil
}

func buildManifestURL(endpoint, repoName, reference string) string {
	return fmt.Sprintf("%s/v2/%s/manifests/%s", endpoint, repoName, reference)
}

func buildTagListURL(endpoint, repoName string) string {
	return fmt.Sprintf("%s/v2/%s/tags/list", endpoint, repoName)
}

func parseError(err error) error {
	if urlErr, ok := err.(*url.Error); ok {
		if regErr, ok := urlErr.Err.(*Error); ok {
			return regErr
		}
	}
	return err
}

// Error wrap HTTP status code and message as an error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error ...
func (e *Error) Error() string {
	return fmt.Sprintf("http error: code %d, message %s", e.Code, e.Message)
}

// String wraps the error msg to the well formatted error message
func (e *Error) String() string {
	data, err := json.Marshal(&e)
	if err != nil {
		return e.Message
	}
	return string(data)
}
