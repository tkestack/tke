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

package smtp

import (
	"testing"
	"time"

	"github.com/golang/glog"
	v1 "tkestack.io/tke/api/notify/v1"
	"tkestack.io/tke/pkg/util/log"
)

func init() {
	logOpts := log.NewOptions()
	logOpts.EnableCaller = true
	logOpts.Level = log.ErrorLevel
	log.Init(logOpts)
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			log.Flush()
		}
	}()
}

func TestSmtpSend(t *testing.T) {
	smtp := &v1.ChannelSMTP{
		SMTPHost: "!!smtp.qq.com",
		SMTPPort: 465,
		Email:    "0000@qq.com",
		Password: "lalnmzdxufjfbiaj",
	}

	template := &v1.TemplateText{
		Body:   "This is body: {{.body}} ",
		Header: "Hi, {{.header}}! ",
	}
	//send to
	email := "464813006@qq.com"

	variables := map[string]string{
		"header": "tony",
		"body":   "this is test",
	}

	header, body, err := Send(smtp, template, email, variables)
	log.Debugf("header: %s", header)
	log.Debugf("body: %s", body)
	log.Debugf("err: %v", err)

	if err != nil {
		glog.Errorf("error %v", err)
		return
	}
}
