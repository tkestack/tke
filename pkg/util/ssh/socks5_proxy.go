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

	"golang.org/x/net/proxy"
	"tkestack.io/tke/pkg/util/log"
)

const defaultCheckTargetAddr = "ccr.ccs.tencentyun.com:443"

type SOCKS5 struct {
	Host            string
	Port            int
	User            string
	Password        string
	DialTimeOut     time.Duration
	CheckTargetAddr string
}

func (sk SOCKS5) ProxyConn(targetAddr string) (net.Conn, func(), error) {
	addr := net.JoinHostPort(sk.Host, fmt.Sprintf("%d", sk.Port))
	dialer, err := proxy.SOCKS5("tcp", addr, nil, proxy.Direct)
	if err != nil {
		log.Errorf("get socks5 dialer to %s failed: %v", addr, err)
		return nil, nil, err
	}

	type result struct {
		conn net.Conn
		err  error
	}
	ch := make(chan result, 1)
	go func() {
		conn, err := dialer.Dial("tcp", targetAddr)
		if err != nil {
			log.Errorf("proxy %s dial %s failed: %v", addr, targetAddr)
		} else {
			log.Debugf("proxy %s dial %s sucess", addr, targetAddr)
		}
		r := result{conn: conn, err: err}
		ch <- r
	}()
	select {
	case r := <-ch:
		return r.conn,
			func() {
				r.conn.Close()
			},
			r.err
	case <-time.After(sk.DialTimeOut):
		return nil, nil, fmt.Errorf("proxy %s dial %s time out in %s", addr, targetAddr, sk.DialTimeOut.String())
	}
}

func (sk SOCKS5) CheckTunnel() error {
	targetAddr := sk.CheckTargetAddr
	if len(targetAddr) == 0 {
		targetAddr = defaultCheckTargetAddr
	}
	_, closer, err := sk.ProxyConn(targetAddr)
	if err != nil {
		return fmt.Errorf("tunnel is unavailable: %v", err)
	}
	defer closer()
	return nil
}
