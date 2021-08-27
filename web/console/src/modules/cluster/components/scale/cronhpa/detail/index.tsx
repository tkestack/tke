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
/**
 * hpa详情页入口
 */
import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { isEmpty, useRefresh } from '@src/modules/common/utils';
import { router } from '@src/modules/cluster/router';
import {
  Layout,
  Card,
  Tabs,
  TabPanel,
  Justify,
  Button,
  Icon
} from '@tencent/tea-component';
import { TipInfo } from '@src/modules/common';
import Basic from './basic';
import Event from './event';
import Yaml from './yaml';

const { Body, Content } = Layout;
const tabs = [
  { id: 'basic', label: '详情' },
  { id: 'event', label: '事件' },
  { id: 'yaml', label: 'YAML' }
];
const Detail = React.memo((props: {
  selectedHpa: any;
}) => {
  const route = useSelector((state) => state.route);
  const urlParams = router.resolve(route);
  const { clusterId } = route.queries;
  const { selectedHpa } = props;
  const { name = '', namespace = '' } = !isEmpty(selectedHpa) ? selectedHpa.metadata : [];
  const { refreshFlag, triggerRefresh } = useRefresh();

  return (
    <Layout>
      <Body>
        <Content>
          <Content.Header
            showBackButton
            onBackButtonClick={() => history.back()}
            title={`${clusterId}/CronHPA:${name}(${namespace})`}
          />
          <Content.Body>
            <Tabs ceiling animated={false} tabs={tabs}>
              <TabPanel id="basic">
                <Card>
                  <Card.Body title={t('基本信息')}>
                    <Basic selectedHpa={selectedHpa} />
                  </Card.Body>
                </Card>
              </TabPanel>
              <TabPanel id="event">
                <TipInfo>{t('资源事件只保存最近1小时内发生的事件，请尽快查阅。')}</TipInfo>
                <Justify
                  style={{ marginBottom: '10px' }}
                  right={
                    <React.Fragment>
                      {/*<span*/}
                      {/*  className="descript-text"*/}
                      {/*  style={{ display: 'inline-block', verticalAlign: 'middle', marginRight: '10px', fontSize: '12px' }}*/}
                      {/*>*/}
                      {/*  {t('自动刷新')}*/}
                      {/*</span>*/}
                      {/*<Text reset>刷新:</Text>*/}
                      <Icon type="refresh" onClick={() => {
                        triggerRefresh();
                      }}
                      />
                      {/*<Switch*/}
                      {/*  onChange={value => {*/}
                      {/*    toggleRefresh(value);*/}
                      {/*  }}*/}
                      {/*  className="mr20"*/}
                      {/*/>*/}
                    </React.Fragment>
                  }
                />
                <Card>
                  <Card.Body>
                    <Event selectedHpa={selectedHpa} refreshFlag={refreshFlag} />
                  </Card.Body>
                </Card>
              </TabPanel>
              <TabPanel id="yaml">
                <Justify
                  style={{ marginBottom: '10px' }}
                  left={
                    <Button
                      type="primary"
                      onClick={() => {
                        router.navigate(
                            { ...urlParams, mode: 'modify-yaml' },
                            route.queries
                        );
                      }}
                    >
                      {t('编辑YAML')}
                    </Button>
                  }
                />
                <Card>
                  <Card.Body >
                    <Yaml />
                  </Card.Body>
                </Card>
              </TabPanel>
            </Tabs>
          </Content.Body>
        </Content>
      </Body>
    </Layout>
  );
});
export default Detail;
