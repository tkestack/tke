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
import { QueryState, RecordSet, uuid } from '@tencent/ff-redux';

import { subRouterConfig } from '../../../config';
import { CLB, SubRouter, SubRouterFilter } from './models';

/** 将各种资源的接口导出 */
export * from './WebAPI/index';

/**
 * subRouter列表的拉取
 * @param query subRouter列表的查询
 */
export async function fetchSubRouterList(query: QueryState<SubRouterFilter>) {
  const { module, clusterId } = query.filter;

  const response = subRouterConfig(module);

  console.log(response);

  const subRouterList: any = await checkRouterVisible(response, { clusterId });

  const result: RecordSet<SubRouter> = {
    recordCount: subRouterList.length,
    records: subRouterList
  };

  return result;
}

async function checkRouterVisible(routeList, { clusterId }) {
  const list = await Promise.all(
    routeList?.map(async ({ visible = _ => Promise.resolve(true), sub, ...other }) => ({
      ...other,

      visible: await visible({ clusterId }),

      sub: sub && (await checkRouterVisible(sub, { clusterId })),

      id: uuid()
    }))
  );

  return list?.filter(({ visible }) => visible) ?? [];
}
