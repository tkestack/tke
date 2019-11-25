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

package wechat

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

func TestWechatSend(t *testing.T) {
	channel := &v1.ChannelWechat{
		AppID:     "!!!wx654580d11ee4650a",
		AppSecret: "a53f88905b9f67647e6faae823fd225b",
	}

	tempalte := &v1.TemplateWechat{
		TemplateID: "DGjbjnYwA9dsclLQ96GLgnC1cSl9FWuYvnJD3UGed8s",
		URL:        "https://baidu.com",
		Body:       "content: {{.body}}",
	}
	openID := "oQ0zg53IUIgQ7xoFbGJGL58ZBVj0"
	variables := map[string]string{
		"body": "success",
	}
	msgID, body, err := Send(channel, tempalte, openID, variables)
	log.Debugf("msgID: %s", msgID)
	log.Debugf("body: %s", body)
	log.Debugf("err: %v", err)
	if err != nil || msgID == "" {
		glog.Error(err)
		return
	}
}
