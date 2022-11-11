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

package webhook

import (
	"net/http"
	"net/url"

	v1 "tkestack.io/tke/api/notify/v1"
	"tkestack.io/tke/pkg/notify/controller/messagerequest/util"
	"tkestack.io/tke/pkg/util/log"
)

// webhookBody represents the body info to request a webhook server
type webhookBody struct {
	Receivers []*v1.Receiver `json:"receivers"`
	Content   string         `json:"content"`
}

// Send notification to webhook server
func Send(channel *v1.ChannelWebhook, template *v1.TemplateText, receivers []*v1.Receiver, variables map[string]string, status string) (content string, err error) {
	content, err = util.ParseTemplate("webhookContent", template.Body, variables)
	if err != nil {
		return "", err
	}

	alertStatus := util.GetAlertStatus(status)
	contentWithAlertStatus := alertStatus + "\r\n" + content //add  alertStatus first
	body := webhookBody{
		Receivers: receivers,
		Content:   contentWithAlertStatus,
	}
	log.Debugf("webhook body: %v", body)
	err = requestToWebhook(channel, body)
	return content, err
}

// requestToWebhook is used to do a post request to webhook server
func requestToWebhook(channel *v1.ChannelWebhook, reqBody interface{}) error {
	reqURL, err := url.Parse(channel.URL)
	if err != nil {
		return err
	}
	option := util.Option{
		Protocol: reqURL.Scheme,
		Host:     reqURL.Host,
		Path:     reqURL.Path,
		Method:   http.MethodPost,
		Body:     reqBody,
		Headers:  map[string]string{"Content-Type": "application/json"},
	}
	_, err = util.Request(option)
	return err
}
