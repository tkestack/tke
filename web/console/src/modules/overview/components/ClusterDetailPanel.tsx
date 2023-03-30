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
import { Bubble, Button, Card, Col, Icon, List, Row, SearchBox, Text } from '@tencent/tea-component';
import * as React from 'react';
import { clusterStatus } from '../constants/Config';
import { ClusterDetail, ClusterOverview } from '../models/RootState';
export function ClusterDetailPanel(props: { clusterData: ClusterOverview }) {
  const { clusterData } = props;
  const isLodingDone = !!clusterData;
  const [search, setSearch] = React.useState('');
  let clusterList = isLodingDone ? clusterData.clusters : [];
  if (search) {
    clusterList = clusterList.filter(item => item.clusterID.includes(search));
  }
  return (
    <Card>
      <Card.Body title="集群状态">
        <SearchBox value={search} onChange={setSearch}></SearchBox>
        {clusterList.map((cluster, index) => _renderClusterCard(cluster, index))}
      </Card.Body>
    </Card>
  );
}

function _renderClusterCard(cluster: ClusterDetail, index: number) {
  const masterErrorTips = [];
  if (!cluster.etcdHealthy) {
    masterErrorTips.push(
      <Text key={'etcd'} style={{ display: 'block' }}>
        etcd异常
      </Text>
    );
  }
  if (!cluster.schedulerHealthy) {
    masterErrorTips.push(
      <Text key={'scheduler'} style={{ display: 'block' }}>
        scheduler异常
      </Text>
    );
  }
  if (!cluster.controllerManagerHealthy) {
    masterErrorTips.push(
      <Text key={'controllerManager'} style={{ display: 'block' }}>
        controllerManager异常
      </Text>
    );
  }

  const isUnlink = cluster.clusterPhase !== 'Running';
  const isFailed = cluster.clusterPhase === 'Failed';

  return (
    <Card style={{ marginTop: index === 0 ? 20 : 1 }} key={cluster.clusterID}>
      <Card.Body
        title={
          <>
            <Button
              disabled={isUnlink}
              type={'link'}
              onClick={() => {
                location.href =
                  location.origin +
                  `/tkestack/cluster/sub/list/resource/deployment?rid=1&clusterId=${cluster.clusterID}&np=default`;
              }}
            >
              {cluster.clusterID}
            </Button>

            <Text theme="weak"> {cluster.clusterDisplayName} </Text>

            {cluster.clusterPhase !== 'Running' && (
              <Bubble content={clusterStatus[cluster.clusterPhase] ? clusterStatus[cluster.clusterPhase].text : '-'}>
                <Icon type={'info'}></Icon>
              </Bubble>
            )}
          </>
        }
        style={{ padding: '15px 10px' }}
      >
        <Row>
          <Col>
            <div style={{ backgroundColor: '#F2F2F2', padding: '16px 10px' }}>
              <Text style={{ fontSize: 18, fontWeight: 500 }}>{isFailed ? '-' : cluster.cpuUsage}</Text>
              <Text style={{ fontSize: 12, fontWeight: 600 }}> CPU利用率</Text>
              <div>
                <Text theme={'label'} reset>{`总数: ${cluster.cpuCapacity} 核 Request已分配: ${
                  isFailed ? '-' : cluster.cpuRequestRate
                }`}</Text>
              </div>
            </div>
          </Col>
          <Col>
            <div style={{ backgroundColor: '#F2F2F2', padding: '16px 10px' }}>
              <Text style={{ fontSize: 18, fontWeight: 500 }}>{isFailed ? '-' : cluster.memUsage}</Text>
              <Text style={{ fontSize: 12, fontWeight: 600 }}> 内存利用率</Text>
              <div>
                <Text theme={'label'} reset>{`总数: ${
                  isFailed ? '-' : (cluster.memCapacity / 1.0 / 1024 / 1024 / 1024).toFixed(3)
                } GB Request已分配: ${isFailed ? '-' : cluster.memRequestRate}`}</Text>
              </div>
            </div>
          </Col>
          <Col>
            <List>
              <List.Item style={{ verticalAlign: 'center' }}>
                节点(
                <Button
                  type={'link'}
                  disabled={isUnlink}
                  onClick={() => {
                    location.href =
                      location.origin +
                      `/tkestack/cluster/sub/list/nodeManage/node?rid=1&clusterId=${cluster.clusterID}&np=default`;
                  }}
                >
                  {isFailed ? '-' : cluster.nodeCount}个
                </Button>
                )
              </List.Item>
              <List.Item style={{ verticalAlign: 'center' }}>
                Workload(
                <Button
                  type={'link'}
                  disabled={isUnlink}
                  onClick={() => {
                    location.href =
                      location.origin +
                      `/tkestack/cluster/sub/list/resource/deployment?rid=1&clusterId=${cluster.clusterID}&np=default`;
                  }}
                >
                  {isFailed ? '-' : cluster.workloadCount}个
                </Button>
                )
              </List.Item>
              <List.Item style={{ verticalAlign: 'center' }}>{'Master&ETCD'}</List.Item>
            </List>
          </Col>
          <Col span={2}>
            <List>
              <List.Item style={{ verticalAlign: 'center' }}>
                {cluster.nodeAbnormal > 0 ? (
                  <>
                    <Text theme={'danger'}>异常</Text>
                    <Bubble content={`该集群中有${cluster.nodeAbnormal}个节点异常`}>
                      <Icon type="info" />
                    </Bubble>
                  </>
                ) : (
                  <Text theme={'success'}>正常</Text>
                )}
              </List.Item>
              <List.Item style={{ verticalAlign: 'center' }}>
                {cluster.workloadAbnormal > 0 ? (
                  <>
                    <Text theme={'danger'}>异常</Text>
                    <Bubble content={`该集群中有${cluster.nodeAbnormal}个工作负载异常`}>
                      <Icon type="info" />
                    </Bubble>
                  </>
                ) : (
                  <Text theme={'success'}>正常</Text>
                )}
              </List.Item>
              <List.Item style={{ verticalAlign: 'center' }}>
                {!cluster.etcdHealthy || !cluster.controllerManagerHealthy || !cluster.schedulerHealthy ? (
                  <>
                    <Text theme={'danger'}>异常</Text>
                    <Bubble content={masterErrorTips}>
                      <Icon type="info" />
                    </Bubble>
                  </>
                ) : (
                  <Text theme={'success'}>正常</Text>
                )}
              </List.Item>
            </List>
          </Col>
        </Row>
      </Card.Body>
    </Card>
  );
}
