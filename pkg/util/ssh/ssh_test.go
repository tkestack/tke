/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package ssh_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"tkestack.io/tke/pkg/util/ssh"

	// env for load env
	_ "tkestack.io/tke/test/util/env"
)

var s *ssh.SSH

func init() {
	port, err := strconv.Atoi(os.Getenv("SSH_PORT"))
	utilruntime.Must(err)
	s, _ = ssh.New(&ssh.Config{
		Host:     os.Getenv("SSH_HOST"),
		Port:     port,
		User:     os.Getenv("SSH_USER"),
		Password: os.Getenv("SSH_PASSWORD"),
	})
}

func TestSudo(t *testing.T) {
	output, err := s.CombinedOutput("whoami")
	assert.Nil(t, err)
	assert.Equal(t, "root", strings.TrimSpace(string(output)))
}

func TestQuote(t *testing.T) {
	output, err := s.CombinedOutput(`echo "a" 'b'`)
	assert.Nil(t, err)
	assert.Equal(t, "a b", strings.TrimSpace(string(output)))
}

func TestWriteFile(t *testing.T) {
	data := []byte("Hello")
	dst := "/tmp/test"

	err := s.WriteFile(bytes.NewBuffer(data), dst)
	assert.Nil(t, err)

	output, err := s.ReadFile(dst)
	assert.Nil(t, err)
	assert.Equal(t, data, output)
}

func TestCoppyFile(t *testing.T) {
	src := os.Args[0]
	srcData, err := ioutil.ReadFile(src)
	assert.Nil(t, err)

	dst := "/tmp/test"
	err = s.CopyFile(src, dst)
	assert.Nil(t, err)

	output, err := s.ReadFile(dst)
	assert.Nil(t, err)

	assert.Equal(t, srcData, output)
}

func TestExist(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"exist",
			args{
				filename: "/tmp",
			},
			true,
			false,
		},
		{
			"not exist",
			args{
				filename: "/tmpfda",
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Exist(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exist() got = %v, want %v", got, tt.want)
			}
		})
	}
}
