/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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

package crypto

import "testing"

func TestAesEncrypt(t *testing.T) {
	type args struct {
		orig string
		key  string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Test AES encrypt and decrypt",
			args{
				orig: "abc",
				key:  NewAesKey(),
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptGot, encryptErr := AesEncrypt(tt.args.orig, tt.args.key)
			decryptGot, decryptErr := AesDecrypt(encryptGot, tt.args.key)
			if decryptGot != tt.args.orig {
				t.Errorf("AesEncrypt() = %v err: %v, AesDecrypt() = %v err: %v, orig is %v, key is %v", encryptGot, encryptErr, decryptGot, decryptErr, tt.args.orig, tt.args.key)
			}
		})
	}
}
