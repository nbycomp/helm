/*
Copyright The Helm Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package getter

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/nbycomp/helm/internal/experimental/registry"
)

// OCIGetter is the default HTTP(/S) backend handler
type OCIGetter struct {
	opts options
}

//Get performs a Get from repo.Getter and returns the body.
func (g *OCIGetter) Get(href string, options ...Option) (*bytes.Buffer, error) {
	for _, opt := range options {
		opt(&g.opts)
	}
	return g.get(href)
}

func (g *OCIGetter) get(href string) (*bytes.Buffer, error) {
	client := g.opts.registryClient

	ref := strings.TrimPrefix(href, "oci://")
	if version := g.opts.version; version != "" {
		ref = fmt.Sprintf("%s:%s", ref, version)
	}

	r, err := registry.ParseReference(ref)
	if err != nil {
		return nil, err
	}

	buf, err := client.PullChart(r)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// NewOCIGetter constructs a valid http/https client as a Getter
func NewOCIGetter(ops ...Option) (Getter, error) {
	registryClient, err := registry.NewClient()
	if err != nil {
		return nil, err
	}

	client := OCIGetter{
		opts: options{
			registryClient: registryClient,
		},
	}

	for _, opt := range ops {
		opt(&client.opts)
	}

	return &client, nil
}
