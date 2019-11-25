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

package installer

import (
	"net"
	"time"

	"tkestack.io/tke/pkg/util/ssh"

	"github.com/emicklei/go-restful"
	"gopkg.in/go-playground/validator.v9"
	"k8s.io/apimachinery/pkg/api/errors"
)

// SSHResource is the REST layer to the Cluster domain
type SSHResource struct {
}

// NewSSHResource create a SSHResource
func NewSSHResource() *SSHResource {
	c := new(SSHResource)

	return c
}

// WebService creates a new service that can handle REST requests for SSHResource.
func (c *SSHResource) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/api/ssh")
	ws.Produces(restful.MIME_JSON)
	ws.Consumes(restful.MIME_JSON)

	ws.Route(ws.POST("").To(c.findSSH).
		Reads(sshConfig{}).Writes(SSH{}))

	return ws
}

type sshConfig struct {
	Host       net.IP `json:"host" validate:"required"`
	Port       int    `json:"port" validate:"required"`
	User       string `json:"user" validate:"required"`
	Password   []byte `json:"password,omitempty"`
	PrivateKey []byte `json:"privateKey,omitempty"`
	PassPhrase []byte `json:"passPhrase,omitempty"`
}

// SSH for test ssh result
type SSH struct {
	Status string `json:"status,omitempty"`
}

func (*SSHResource) findSSH(request *restful.Request, response *restful.Response) {
	req := new(sshConfig)

	apiStatus := func() errors.APIStatus {
		err := request.ReadEntity(req)
		if err != nil {
			return errors.NewBadRequest(err.Error())
		}

		validate := validator.New()
		err = validate.Struct(req)
		if err != nil {
			return errors.NewBadRequest(err.Error())
		}

		sshConfig := &ssh.Config{
			User:        req.User,
			Host:        req.Host.String(),
			Port:        req.Port,
			Password:    string(req.Password),
			PrivateKey:  req.PrivateKey,
			PassPhrase:  req.PassPhrase,
			DialTimeOut: time.Second,
			Retry:       0,
		}
		s, err := ssh.New(sshConfig)
		if err != nil {
			return errors.NewBadRequest(err.Error())
		}

		err = s.Ping()
		if err != nil {
			return errors.NewBadRequest(err.Error())
		}

		return nil
	}()

	if apiStatus != nil {
		response.WriteHeaderAndJson(int(apiStatus.Status().Code), apiStatus.Status(), restful.MIME_JSON)
	} else {
		response.WriteAsJson(SSH{Status: "OK"})
	}
}
