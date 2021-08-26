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

/**
 * hpa详情，基本信息tab
 */
import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { isEmpty } from '@src/modules/common/utils';
import { router } from '@src/modules/cluster/router';
import {
  Layout,
  Text,
  Bubble,
  Icon
} from '@tencent/tea-component';
import { MetricsResourceMap } from '../constant';

const Basic = React.memo((props: {
  selectedHpa: any;
}) => {
  // const route = useSelector((state) => state.route);
  // const urlParams = router.resolve(route);
  const { selectedHpa } = props;

  return !isEmpty(selectedHpa) && (
    <ul className="item-descr-list">
      <li>
        <span className="item-descr-tit"><Trans>名称</Trans></span>
        <span className="item-descr-txt">{selectedHpa.metadata.name}</span>
      </li>
      <li>
        <span className="item-descr-tit">Namespace</span>
        <span className="item-descr-txt">{selectedHpa.metadata.namespace}</span>
      </li>
      <li>
        <span className="item-descr-tit"><Trans>关联工作负载</Trans></span>
        <span className="item-descr-txt">{selectedHpa.spec.scaleTargetRef.kind}:{selectedHpa.spec.scaleTargetRef.name}</span>
      </li>
      <li>
        <span className="item-descr-tit">
          <Trans>CronHPA指标</Trans>
          <Bubble content={t('根据设置的Crontab（Crontab语法格式，例如 "0 23 * * 5"表示每周五23:00）周期性地设置实例数量')}>
            <Icon
              type="info"
              style={{ marginLeft: '5px', cursor: 'pointer', verticalAlign: 'text-bottom' }}
            />
          </Bubble>
        </span>
        <span className="item-descr-txt">
          {selectedHpa.spec.crons.map((item, index) => {
            const { schedule, targetReplicas } = item;
            return <Text key={index} parent="div">{`${schedule} ${targetReplicas}个`}</Text>;
          })}
        </span>
      </li>
    </ul>
  );
});
export default Basic;
