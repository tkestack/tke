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

package util

func ReadSystemdUnitLog(unitName string, line uint64) ([]byte, error) {
	//r, err := sdjournal.NewJournalReader(sdjournal.JournalReaderConfig{
	//	NumFromTail: line,
	//	Matches: []sdjournal.Match{
	//		{
	//			Field: "_SYSTEMD_UNIT",
	//			Value: fmt.Sprintf(`%s`, unitName),
	//		},
	//	},
	//	Formatter: func(entry *sdjournal.JournalEntry) (string, error) {
	//		return fmt.Sprintf("%s\n", entry.Fields["MESSAGE"]), nil
	//	},
	//})
	//if err != nil {
	//	return []byte{}, nil
	//}
	//
	//if r == nil {
	//	return []byte{}, errors.New("")
	//}
	//
	//defer r.Close()
	//
	//buf := new(bytes.Buffer)
	//var e error
	//for c := -1; c != 0 && e == nil; {
	//	tmpBuf := make([]byte, 5)
	//	c, e = r.Read(tmpBuf)
	//	if c > len(tmpBuf) {
	//		return []byte{}, errors.New(fmt.Sprintf("Got unexpected read length: %d vs %d", c, len(tmpBuf)))
	//	}
	//	_, _ = buf.Write(tmpBuf)
	//}

	return []byte{}, nil
}
