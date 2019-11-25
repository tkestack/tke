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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/golang/glog"
	v1 "tkestack.io/tke/api/notify/v1"
	"tkestack.io/tke/pkg/notify/controller/messagerequest/util"
	"tkestack.io/tke/pkg/util/log"
)

const (
	getTokenURL          = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential"
	wechatSendMessageURL = "https://api.weixin.qq.com/cgi-bin/message/template/send"
)

type miniProgramInfo struct {
	MiniProgramAppID    string `json:"appid"`
	MiniProgramPagePath string `json:"pagepath,omitempty"`
}

type dataItem struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

type bodyInfo struct {
	Touser      string              `json:"touser"`
	TemplateID  string              `json:"template_id"`
	URL         string              `json:"url,omitempty"`
	MiniProgram miniProgramInfo     `json:"miniprogram,omitempty"`
	Data        map[string]dataItem `json:"data"`
}

var resTokenBody struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

type resMessageBody struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	MsgID   int    `json:"msgid"`
}

type item struct {
	token           string
	expiresIn       int64
	accessTimestamp int64
}

var tokenCache sync.Map

// Send notification by wechat
func Send(channel *v1.ChannelWechat, template *v1.TemplateWechat, openID string, variables map[string]string) (messageID string, body string, err error) {
	body, err = util.ParseTemplate("wechatBody", template.Body, variables)
	if err != nil {
		return "", "", err
	}

	if util.SelfdefineURL != "" {
		selfdefineReqBody := util.SelfdefineBodyInfo{
			Type: "wechat",
			Body: body,
		}
		err = util.RequestToSelfdefine(selfdefineReqBody)
		return "", body, err
	}

	accessToken, err := getToken(channel, openID)
	if err != nil {
		return "", body, err
	}
	reqURL, err := url.Parse(wechatSendMessageURL)
	if err != nil {
		return "", body, err
	}

	miniProgram := miniProgramInfo{
		MiniProgramAppID:    template.MiniProgramAppID,
		MiniProgramPagePath: template.MiniProgramPagePath,
	}

	reqBody := bodyInfo{
		Touser:      openID,
		TemplateID:  template.TemplateID,
		URL:         template.URL,
		MiniProgram: miniProgram,
		Data:        mapToDataItemMap(variables),
	}

	option := util.Option{
		Protocol: reqURL.Scheme,
		Host:     reqURL.Host,
		Path:     reqURL.Path + "?access_token=" + accessToken,
		Method:   http.MethodPost,
		Headers:  map[string]string{"Content-Type": "application/json"},
		Body:     reqBody,
	}

	var resMessage resMessageBody
	response, err := util.Request(option)
	if err != nil {
		log.Errorf("Request error %v", err)
		return "", body, err
	}
	err = json.Unmarshal(response, &resMessage)
	if err != nil {
		return "", body, err
	}

	if resMessage.ErrCode != 0 {
		err = fmt.Errorf("post wechatSendMessageURL error: errcode=%v, errmsg=%v", resMessage.ErrCode, resMessage.ErrMsg)
		return "", body, err
	}

	return strconv.Itoa(resMessage.MsgID), body, nil
}

func getToken(channel *v1.ChannelWechat, openID string) (string, error) {
	var accessToken string
	if cacheItem, ok := tokenCache.Load(openID); !ok || util.GetCurrentTime()-cacheItem.(*item).accessTimestamp >= cacheItem.(*item).expiresIn {
		reqURL := getTokenURL + "&appid=" + channel.AppID + "&secret=" + channel.AppSecret
		glog.Info("get access_token via post request")
		resp, err := http.Get(reqURL)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&resTokenBody)
		if err != nil {
			return "", err
		}
		if resTokenBody.ErrCode != 0 {
			err = fmt.Errorf("get access_token error: errcode=%v, errmsg=%v", resTokenBody.ErrCode, resTokenBody.ErrMsg)
			return "", err
		}
		accessToken = resTokenBody.AccessToken

		newItem := &item{
			token:           accessToken,
			expiresIn:       resTokenBody.ExpiresIn,
			accessTimestamp: util.GetCurrentTime(),
		}
		tokenCache.Store(openID, newItem)
	} else {
		glog.Info("get access_token from cache")
		accessToken = cacheItem.(*item).token
	}

	return accessToken, nil
}

func mapToDataItemMap(m map[string]string) map[string]dataItem {
	s := make(map[string]dataItem)
	for k, v := range m {
		s[k] = dataItem{
			Value: v,
		}
	}
	return s
}
