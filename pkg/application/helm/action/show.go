/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package action

import (
	"os"
	"strings"

	securejoin "github.com/cyphar/filepath-securejoin"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"tkestack.io/tke/pkg/util/file"
)

// ShowOptions is the installation options to a install call.
type ShowOptions struct {
	ChartPathOptions
}

// ShowInfo is the result of a chart show call.
type ShowInfo struct {
	Values map[string]string
	Readme map[string]string
	Chart  *chart.Chart
}

// Show installs a chart archive
func (c *Client) Show(options *ShowOptions) (show ShowInfo, err error) {
	client := NewShow(action.ShowAll)

	options.ChartPathOptions.ApplyTo(&client.ChartPathOptions)

	settings, err := NewSettings(options.ChartRepo)
	if err != nil {
		return show, err
	}

	// unpack first if need
	root := settings.RepositoryCache
	if options.ExistedFile != "" && file.IsFile(options.ExistedFile) {
		temp, err := ExpandFile(options.ExistedFile, settings.RepositoryCache)
		if err != nil {
			return show, err
		}
		root = temp
		defer func() {
			os.RemoveAll(temp)
		}()
	}
	chartDir, err := securejoin.SecureJoin(root, options.Chart)
	if err != nil {
		return show, err
	}

	cp, err := client.ChartPathOptions.LocateChart(chartDir, settings)
	if err != nil {
		return show, err
	}

	client.OutputFormat = action.ShowValues
	values, err := client.Run(cp)
	if err != nil {
		return show, err
	}
	client.OutputFormat = action.ShowReadme
	readme, err := client.Run(cp)
	if err != nil {
		return show, err
	}
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return show, err
	}

	return ShowInfo{Values: values, Readme: readme, Chart: chartRequested}, nil
}

var readmeFileNames = []string{"readme.md", "readme.txt", "readme"}

// Show is the action for checking a given release's information.
//
// It provides the implementation of 'helm show' and its respective subcommands.
type Show struct {
	action.Show
	chart *chart.Chart // for testing
}

// NewShow creates a new Show object with the given configuration.
func NewShow(output action.ShowOutputFormat) *Show {
	s := &Show{}
	s.Show = action.Show{
		OutputFormat: output,
	}
	return s
}

// Run executes 'helm show' against the given release.
func (s *Show) Run(chartpath string) (map[string]string, error) {
	if s.chart == nil {
		chrt, err := loader.Load(chartpath)
		if err != nil {
			return nil, err
		}
		s.chart = chrt
	}

	out := make(map[string]string)

	if (s.OutputFormat == action.ShowValues) && s.chart.Values != nil {
		for _, f := range s.chart.Raw {
			if f.Name == chartutil.ValuesfileName {
				out[f.Name] = string(f.Data)
			}
		}
	}

	if s.OutputFormat == action.ShowReadme {
		readme := findReadme(s.chart.Files)
		if readme == nil {
			return out, nil
		}
		out[readme.Name] = string(readme.Data)
	}
	return out, nil
}

func findReadme(files []*chart.File) (file *chart.File) {
	for _, file := range files {
		for _, n := range readmeFileNames {
			if strings.EqualFold(file.Name, n) {
				return file
			}
		}
	}
	return nil
}
