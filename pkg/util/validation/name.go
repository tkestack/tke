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

package validation

import (
	"fmt"
	"regexp"
)

const dns1123NameMaxLength int = 32
const dns1123NameMinLength int = 3
const displayNameMaxLength int = 255

var dns1123NameFmt = "^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"
var emailFmt = `^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,4}$`
var phoneNumberFmt = `^1[3|4|5|7|8][0-9]{9}$`

var dns1123NameRegexp = regexp.MustCompile(dns1123NameFmt)
var emailRegexp = regexp.MustCompile(emailFmt)
var phoneNumberRegexp = regexp.MustCompile(phoneNumberFmt)

// IsDNS1123Name tests for a string that conforms to the definition of a name in
// DNS (RFC 1123).
func IsDNS1123Name(value string) error {
	if value == "" {
		return fmt.Errorf("must be specified")
	}
	if len(value) < dns1123NameMinLength || len(value) > dns1123NameMaxLength {
		return fmt.Errorf("length must be greater than %d and less than %d", dns1123NameMinLength, dns1123NameMaxLength)
	}
	if !dns1123NameRegexp.MatchString(value) {
		return fmt.Errorf("must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character")
	}
	return nil
}

// IsDisplayName test whether the given value meets the specification of the
// display name.
func IsDisplayName(value string) error {
	if value == "" {
		return fmt.Errorf("must be specified")
	}
	if len(value) > displayNameMaxLength {
		return fmt.Errorf("length must be less than %d", displayNameMaxLength)
	}
	return nil
}

// IsEmail test whether the given value meets the specification of the email.
func IsEmail(value string) error {
	if value == "" {
		return fmt.Errorf("must be specified")
	}

	if !emailRegexp.MatchString(value) {
		return fmt.Errorf("email is not valid format, must satisfy regex %s, examples: 123@abc.com ", emailRegexp)
	}

	return nil
}

// IsPhoneNumber test whether the given value meets the specification of the phone number.
func IsPhoneNumber(value string) error {
	if value == "" {
		return fmt.Errorf("must be specified")
	}

	if !phoneNumberRegexp.MatchString(value) {
		return fmt.Errorf("phoneNumer is not valid format, must satisfy regex %s, examples: 13611111111", emailRegexp)
	}

	return nil
}
