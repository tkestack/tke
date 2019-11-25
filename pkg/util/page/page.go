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

package page

import (
	"strconv"

	"github.com/emicklei/go-restful"
)

const (
	defaultPage = 0
	defaultSize = 0
)

// Pagination contains data for ping items
type Pagination struct {
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	Total      int         `json:"total"`
	TotalPages int         `json:"totalPages"`
	Items      interface{} `json:"items"`
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Pagein returns start and end index by given param.
func Pagein(page int, size int, total int) (start int, end int, pagin Pagination) {
	var (
		totalPages int
	)

	if size <= 0 || page < 0 {
		start = 0
		end = total
		totalPages = 1
	} else {
		start = min(page*size, total)
		end = min(start+size, total)
		totalPages = (total + size - 1) / size
	}

	pagin = Pagination{
		Page:       page,
		PageSize:   size,
		Total:      total,
		TotalPages: totalPages,
	}

	return
}

// ParsePageParam parses page and size from http request.
func ParsePageParam(r *restful.Request) (page int, size int) {
	page, err := strconv.Atoi(r.QueryParameter("page"))
	if err != nil {
		page = defaultPage
	}
	size, err = strconv.Atoi(r.QueryParameter("page_size"))
	if err != nil {
		size = defaultSize
	}
	return
}
