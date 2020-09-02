/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

package version

import (
	"fmt"
	"strconv"
	"strings"
)

// You can parse version like：
//  "1.0";
//  "1.0.1.20140402";
//  "2.0.1.-rc1";
//  "2.11.1.20140402a1";
//  "1.0.0+build1"
//  "1.0build1.alpha2"

const (
	// MaxLen is the max length of version string
	MaxLen = 100
)

// Parse passe version string and return string array. e.g.
//  1.0.1.201.build1
// will return
//  []string{"1", "0", "1", "201", "build", "1"}
// The connector between version numbers can be "-"; "."; "+"; " "
func Parse(str string) (ret []string, err error) {
	l := len(str)
	if l > MaxLen {
		return nil, fmt.Errorf("The maximum length must not be greater than [%v]", MaxLen)
	}

	var start, index int
	var preIsAlpha = false // if last char is alpha([a-zA-Z])
	var v rune

	// currIsAlpha: if current char is alpha
	getRet := func(currIsAlpha bool) {
		if preIsAlpha == currIsAlpha {
			return
		}

		preIsAlpha = currIsAlpha
		if start == index { // skip blank value
			return
		}

		ret = append(ret, str[start:index])
		start = index
	}

	for index, v = range str {
		switch {
		case v == '.' || v == '-' || v == '+' || v == ' ': // special
			if start == index { // skip blank value
				start = index + 1
				continue
			}
			ret = append(ret, str[start:index])
			start = index + 1 // skip current char
		case v >= 48 && v <= 59: // number
			getRet(false)
		case (v >= 65 && v <= 90) || (v >= 97 && v <= 122): // alpha
			getRet(true)
		default:
			return nil, fmt.Errorf("Can not parse [%v]", str[index:index+1])
		}
	}

	if start < l { // last str
		ret = append(ret, str[start:l])
	}

	return
}

// CompareFunc use custom func to compare two versions
func CompareFunc(v1, v2 string, comp func(word1, word2 string) int) int {
	if comp == nil {
		return Compare(v1, v2)
	}

	vv1, err := Parse(v1)
	if err != nil {
		panic(err)
	}
	vv2, err := Parse(v2)
	if err != nil {
		panic(err)
	}

	l2 := len(vv2)
	for index, word1 := range vv1 {
		if l2 <= index {
			return comp(word1, "")
		}

		if ret := comp(word1, vv2[index]); ret == 0 {
			continue
		} else {
			return ret
		}
	}

	l1 := len(vv1)
	if l2 > l1 {
		return comp("", vv2[l1])
	}

	return 0
}

// Compare compare v1 and v2, and return
// >0: v1 has higher version
// =0: equal
// <0: v2 has higher version
func Compare(v1, v2 string) int {
	return CompareFunc(v1, v2, defaultCompare)
}

// version suffix word, later version is bigger
const (
	unknown = iota
	alpha
	beta
	rc
	rtm
	build
	none // keep in the final
)

var suffix = map[string]int{
	"":      none,
	"build": build,
	"rtm":   rtm,
	"rc":    rc,
	"beta":  beta,
	"b":     beta,
	"alpha": alpha,
	"a":     alpha,
}

// change a word to a number
// state： 0 means empty, 1 means suffix word, 2 means normal conversion
func atoi(word string) (num int, state int) {
	var found bool

	switch {
	case len(word) == 0:
		num, state = 0, 0
	case word[0] > 59: // suffix word
		state = 1
		if num, found = suffix[strings.ToLower(word)]; !found {
			num = unknown
		}
	default:
		state = 2
		num1, err := strconv.Atoi(word)
		if err != nil {
			panic(err)
		}
		num = num1
	}
	return
}

// m is a comparison table
//
//  switch v1State {
//	case 0: // empty
//		switch v2State {
//		case 0:
//			return 0
//		case 1:
//			return 1
//		case 2:
//			return -1
//		}
//	case 1: // suffix
//		switch v2State {
//		case 0:
//			return -1
//		case 1:
//			return v1 - v2
//		case 2:
//			return -1
//		}
//	case 2: // normal conversion
//		switch v2State {
//		case 0:
//			return 1
//		case 1:
//			return 1
//		case 2:
//			return v1 - v2
//		}
//	}
var m = [][]int{
	{0, 1, -1},
	{-1, 2, -1},
	{1, 1, 2},
}

func defaultCompare(word1, word2 string) int {
	v1, v1State := atoi(word1)
	v2, v2State := atoi(word2)

	v := m[v1State][v2State]
	if v == 2 {
		return v1 - v2
	}
	return v
}
