import * as React from 'react';
import { Card, Row, Col, MetricsBoard, Icon, Text, Bubble } from '@tencent/tea-component';
import { ClusterOverview } from '../models/RootState';
export function ClusterOverviewPanel(props: { clusterData: ClusterOverview }) {
  let { clusterData } = props;
  let isLodingDone = !!clusterData;

  function projectRender(count) {
    if (count === -1) {
      return (
        <Bubble content="获取业务数量失败">
          <Text>?</Text>
        </Bubble>
      );
    } else if (count === -2) {
      return (
        <Bubble content="集群未安装业务模块">
          <Text>-</Text>
        </Bubble>
      );
    } else {
      return (
        <>
          <Text>{clusterData.projectCount}</Text>
          <Text style={{ fontSize: 14 }}>{'个'}</Text>
        </>
      );
    }
  }

  return (
    <Card>
      <Card.Body title="资源概览" style={{ paddingBottom: 48 }}>
        <Row showSplitLine>
          <Col>
            <MetricsBoard
              title="集群"
              value={
                isLodingDone ? (
                  <>
                    <Text>{clusterData.clusterCount}</Text>
                    <Text style={{ fontSize: 14 }}>{'个'}</Text>
                    {clusterData.clusterAbnormal > 0 && (
                      <Text theme="danger" style={{ fontSize: 14 }}>
                        （异常 {clusterData.clusterAbnormal} 个）
                      </Text>
                    )}
                  </>
                ) : (
                  '-'
                )
              }
            />
          </Col>
          <Col>
            <MetricsBoard
              title="节点"
              value={
                isLodingDone ? (
                  <>
                    <Text>{clusterData.nodeCount}</Text>
                    <Text style={{ fontSize: 14 }}>{'个'}</Text>
                    {clusterData.nodeAbnormal > 0 && (
                      <Text theme="danger" style={{ fontSize: 14 }}>
                        （异常 {clusterData.nodeAbnormal} 个）
                      </Text>
                    )}
                  </>
                ) : (
                  '-'
                )
              }
            />
          </Col>
          <Col>
            <MetricsBoard
              title="负载数"
              value={
                isLodingDone ? (
                  <>
                    <Text>{clusterData.workloadCount}</Text>
                    <Text style={{ fontSize: 14 }}>{'个'}</Text>
                    {clusterData.workloadAbnormal > 0 && (
                      <Text theme="danger" style={{ fontSize: 14 }}>
                        （异常 {clusterData.workloadAbnormal} 个）
                      </Text>
                    )}
                  </>
                ) : (
                  '-'
                )
              }
            />
          </Col>
          <Col>
            <MetricsBoard
              title="业务"
              value={
                isLodingDone ? (
                  <>
                    {projectRender(clusterData.projectCount)}
                    {clusterData.projectAbnormal > 0 && (
                      <Text theme="danger" style={{ fontSize: 14 }}>
                        （异常 {clusterData.projectAbnormal} 个）
                      </Text>
                    )}
                  </>
                ) : (
                  '-'
                )
              }
            />
          </Col>
        </Row>
      </Card.Body>
    </Card>
  );
}
