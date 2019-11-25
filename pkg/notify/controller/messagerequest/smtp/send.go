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
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"

	v1 "tkestack.io/tke/api/notify/v1"
	"tkestack.io/tke/pkg/notify/controller/messagerequest/util"
	"tkestack.io/tke/pkg/util/log"
)

// Send notification by smtp
func Send(channel *v1.ChannelSMTP, template *v1.TemplateText, email string, variables map[string]string) (header, body string, err error) {
	// 这里的header和body分别表示模版渲染后的邮件标题和邮件内容
	header, err = util.ParseTemplate("smtpHeader", template.Header, variables)
	if err != nil {
		return "", "", err
	}
	log.Debugf("header: %v", header)

	body, err = util.ParseTemplate("smtpBody", template.Body, variables)
	if err != nil {
		return header, "", err
	}
	log.Debugf("body: %v", body)

	if util.SelfdefineURL != "" {
		selfdefineReqBody := util.SelfdefineBodyInfo{
			Type:   "smtp",
			Header: header,
			Body:   body,
		}
		err = util.RequestToSelfdefine(selfdefineReqBody)
		return header, body, err
	}

	auth := smtp.PlainAuth("", channel.Email, channel.Password, channel.SMTPHost)

	headers := make(map[string]string)
	headers["From"] = channel.Email
	headers["To"] = email
	headers["Subject"] = header
	headers["Content-Type"] = "text/plain; charset=UTF-8"

	msg := ""
	for k, v := range headers {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + body

	log.Debug("msg", log.String("msg", msg[:]))

	err = sendMail(fmt.Sprintf("%s:%d", channel.SMTPHost, channel.SMTPPort),
		auth, channel.Email, []string{email}, []byte(msg))
	if err != nil {
		log.Errorf("sendMail error: %v", err)
	}
	return header, body, err
}

func sendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("invalid address: %s", err)
	}
	tlsConfig := &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: true,
	}
	var c *smtp.Client

	if port == "465" {
		//via TLS
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return err
		}
		c, err = smtp.NewClient(conn, host)
		if err != nil {
			return err
		}
	} else {
		c, err = smtp.Dial(addr)
		if err != nil {
			return err
		}
	}
	defer c.Close()

	if err := c.Hello("localhost"); err != nil {
		return err
	}

	if port != "465" {
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err := c.StartTLS(tlsConfig); err != nil {
				return err
			}
		}
	}

	if ok, _ := c.Extension("AUTH"); ok {
		if err = c.Auth(a); err != nil {
			return err
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}
