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
export const formatRequestRequest = (rsp: any) => {
  let dataList = [],
    data = null,
    length = 0,
    isAuthorized = true,
    isLoginedSec = true,
    message = '',
    redirect = '';
  if (rsp.data.code === 0) {
    if (rsp.data.data.list) {
      dataList = rsp.data.data.list;
      length = rsp.data.data.total;
    } else {
      data = rsp.data.data;
    }
  } else if (rsp.data.code === 1900) {
    isAuthorized = false;
    isLoginedSec = false;
    message = rsp.data.message;
    redirect = rsp.data.redirect;
  } else if (rsp.data.code === 1800) {
    isAuthorized = false;
    message = rsp.data.message;
    redirect = rsp.data.redirect;
  }

  return {
    data,
    dataList,
    length,
    isAuthorized,
    isLoginedSec,
    message,
    redirect
  };
};
