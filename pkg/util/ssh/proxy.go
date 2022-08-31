/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2022 Tencent. All Rights Reserved.
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

package ssh

import (
	"fmt"
	"net"
	"time"

	"tkestack.io/tke/pkg/util/log"
)

type Proxy interface {
	ProxyConn(targetAddr string) (conn net.Conn, closer func(), err error)
	CheckTunnel() (err error)
}

var _ Proxy = JumpServer{}

type JumpServer struct {
	Config
}

func (sj JumpServer) ProxyConn(targetAddr string) (net.Conn, func(), error) {
	sshstruct, err := New(&sj.Config)
	if err != nil {
		return nil, nil, err
	}
	// do not use sudo in jump server
	sshstruct.Sudo = false
	if sshstruct.DialTimeOut == 0 {
		sshstruct.DialTimeOut = time.Second
	}
	jumperClient, closer, err := sshstruct.newClient()
	if err != nil {
		return nil, nil, err
	}
	type result struct {
		conn net.Conn
		err  error
	}
	ch := make(chan result, 1)
	go func() {
		conn, err := jumperClient.Dial("tcp", targetAddr)
		if err != nil {
			closer()
			log.Errorf("proxy %s dial %s failed: %v", sj.Host, targetAddr)
		} else {
			log.Debugf("proxy %s dial %s sucess", sj.Host, targetAddr)
		}
		r := result{conn: conn, err: err}
		ch <- r
	}()
	select {
	case r := <-ch:
		return r.conn,
			func() {
				closer()
			},
			r.err
	case <-time.After(sshstruct.DialTimeOut):
		return nil, nil, fmt.Errorf("proxy %s dial %s time out in %s", sshstruct.Host, targetAddr, sshstruct.DialTimeOut.String())
	}
}

func (sj JumpServer) CheckTunnel() error {
	sshstruct, err := New(&sj.Config)
	if err != nil {
		return err
	}
	if sshstruct.DialTimeOut == 0 {
		sshstruct.DialTimeOut = time.Second
	}
	_, closer, err := sshstruct.newClient()
	if err != nil {
		return err
	}
	closer()
	return nil
}
