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

package model

// SavedResponse defines the return structure of the chartmuseum server after
// the upload operation is successful.
type SavedResponse struct {
	Saved bool `json:"saved"`
}

// DeletedResponse defines the return structure of the chartmuseum server after
// the delete chart version operation is successful.
type DeletedResponse struct {
	Deleted bool `json:"deleted"`
}
