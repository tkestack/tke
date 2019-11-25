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

package tencentcloudsms

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	v1 "tkestack.io/tke/api/notify/v1"
	"tkestack.io/tke/pkg/notify/controller/messagerequest/util"
	"tkestack.io/tke/pkg/util/log"
)

const (
	chineseCode        = "86"
	tencentCloudSmsURL = "https://yun.tim.qq.com/v5/tlssmssvr/sendsms"
)

type telInfo struct {
	NationCode string `json:"nationcode"`
	Mobile     string `json:"mobile"`
}

type bodyInfo struct {
	Tel    telInfo  `json:"tel"`
	Sign   string   `json:"sign,omitempty"`
	TplID  int      `json:"tpl_id"`
	Params []string `json:"params"`
	Sig    string   `json:"sig"`
	Time   int64    `json:"time"`
	Extend string   `json:"extend,omitempty"`
	Ext    string   `json:"ext,omitempty"`
}

type resBody struct {
	Result int    `json:"result"`
	Errmsg string `json:"errmsg"`
	Ext    string `json:"ext,omitempty"`
	Fee    int    `json:"fee"`
	Sid    string `json:"sid,omitempty"`
}

// Send notification by tencent cloud sms gateway.
func Send(channel *v1.ChannelTencentCloudSMS, template *v1.TemplateTencentCloudSMS, mobile string, variables map[string]string) (messageID string, body string, err error) {
	body, err = util.ParseTemplate("smsBody", template.Body, variables)
	if err != nil {
		return "", "", err
	}

	if util.SelfdefineURL != "" {
		selfdefineReqBody := util.SelfdefineBodyInfo{
			Type: "sms",
			Body: body,
		}
		err = util.RequestToSelfdefine(selfdefineReqBody)
		return "", body, err
	}

	reqURL, err := url.Parse(tencentCloudSmsURL)
	if err != nil {
		return "", body, err
	}

	random := getRandom()
	now := util.GetCurrentTime()

	tel := telInfo{
		NationCode: chineseCode,
		Mobile:     mobile,
	}
	tplID, err := strconv.Atoi(template.TemplateID)
	if err != nil {
		return "", body, err
	}

	reqBody := bodyInfo{
		Tel:    tel,
		Sign:   template.Sign,
		TplID:  tplID,
		Params: mapToSlice(variables, template.Body),
		Sig:    calculateSignature(channel.AppKey, random, now, mobile),
		Time:   now,
	}

	option := util.Option{
		Protocol: reqURL.Scheme,
		Host:     reqURL.Host,
		Path:     reqURL.Path + "?sdkappid=" + channel.SdkAppID + "&random=" + strconv.Itoa(random),
		Method:   http.MethodPost,
		Headers:  map[string]string{"Content-Type": "application/json"},
		Body:     reqBody,
	}

	var resMessage resBody
	response, err := util.Request(option)
	if err != nil {
		log.Errorf("Request error: %v", err)
		return "", body, err
	}
	err = json.Unmarshal(response, &resMessage)
	if err != nil {
		return "", body, err
	}

	if resMessage.Result != 0 {
		err = fmt.Errorf("post tencentCloudSmsURL error: errcode=%v, errmsg=%v", resMessage.Result, resMessage.Errmsg)
		return "", body, err
	}

	return resMessage.Sid, body, nil
}

func getRandom() int {
	min := 100000
	max := 999999
	return rand.Intn(max-min) + min
}

func calculateSignature(appKey string, random int, time int64, phoneNumber string) string {
	var err error
	h := sha256.New()

	_, err = h.Write([]byte("appkey=" + appKey + "&random=" + strconv.Itoa(random) +
		"&time=" + strconv.FormatInt(time, 10) + "&mobile=" + phoneNumber))
	if err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}

func mapToSlice(variables map[string]string, templateBody string) []string {
	re := regexp.MustCompile(`{{\.([a-zA-Z_][a-zA-Z0-9_]*)}}`)
	keys := re.FindAllSubmatch([]byte(templateBody), -1)
	s := make([]string, 0, len(keys))
	for _, key := range keys {
		value, ok := variables[string(key[1])]
		if ok {
			s = append(s, value)
		} else {
			s = append(s, "")
		}
	}
	return s
}
