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

func TestSmsSend(t *testing.T) {
	channel := &v1.ChannelTencentCloudSMS{
		AppKey:   "!!5ca5807e3cb8842523c2437633b210ed",
		SdkAppID: "1400055053",
	}
	template := &v1.TemplateTencentCloudSMS{
		TemplateID: "372399",
		Sign:       "",
		Body:       "alertname:{{.alertName}}, instance:{{.instance}}, startTime:{{.startsAt}}",
	}
	variables := map[string]string{
		"body":      "testContent",
		"alertName": "pod",
		"startsAt":  "2019-09-15",
	}
	testMobile := "15600564755"
	msgID, body, err := Send(channel, template, testMobile, variables)

	log.Debugf("msgID: %s", msgID)
	log.Debugf("body: %s", body)
	log.Debugf("err: %v", err)
	if err != nil || msgID == "" {
		glog.Errorf("send sms error: %v", err)
		return
	}
}

func TestCalSig(t *testing.T) {
	strMobile := "13788888888"
	strAppKey := "5f03a35d00ee52a21327ab048186a2c4"
	strRand := 7226249334
	strTime := int64(1457336869)
	calSig := calculateSignature(strAppKey, strRand, strTime, strMobile)
	expectedSig := "ecab4881ee80ad3d76bb1da68387428ca752eb885e52621a3129dcf4d9bc4fd4"
	if calSig != expectedSig {
		t.Fatalf("calSig is error %s", calSig)
	}
}

func TestMapToSlice(t *testing.T) {
	templateBody := "alertname:{{.name}}, startTime:{{.time}}, content{{.type}}"
	variables := map[string]string{
		"body":   "test_content",
		"time":   "2019",
		"name":   "aaa",
		"unused": "unused",
	}
	expectedSlice := []string{variables["name"], variables["time"], ""}
	res := mapToSlice(variables, templateBody)
	if len(expectedSlice) != len(res) {
		t.Fatal("mapToSlice is error")
	}
	for i := range expectedSlice {
		if res[i] != expectedSlice[i] {
			t.Fatal("mapToSlice is error")
		}
	}
}
