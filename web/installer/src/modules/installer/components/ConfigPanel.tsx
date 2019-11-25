import * as React from "react";
import { RootProps } from "./InstallerApp";
import { Button, Form, Card } from "@tencent/tea-component";

export class ConfigPanel extends React.Component<RootProps> {
  render() {
    const { licenseConfig, actions } = this.props;
    return licenseConfig.cluster ? (
      <div style={{ maxWidth: "1000px", minHeight: "600px", margin: "0 auto" }}>
        <h2 style={{ margin: "40px 0px", fontWeight: 600 }}>TKE Enterprise</h2>
        <Card className="tc-panel">
          <Card.Body>
            <div style={{ padding: "60px 60px 20px", fontSize: "14px" }}>
              <h2 style={{ marginBottom: "20px" }}>License 授权信息</h2>
              <Form.Title>基本信息</Form.Title>
              <Form layout="fixed">
                <Form.Item label="名称">
                  <Form.Text>{licenseConfig.name}</Form.Text>
                </Form.Item>
                <Form.Item label="生效时间">
                  <Form.Text>{licenseConfig.beginDate}</Form.Text>
                </Form.Item>
                <Form.Item label="失效时间">
                  <Form.Text>{licenseConfig.expireDate}</Form.Text>
                </Form.Item>
              </Form>
              <hr />
              <Form.Title>集群限制</Form.Title>
              <Form layout="fixed">
                <Form.Item label="最大集群数">
                  <Form.Text>{licenseConfig.cluster.maxNum} 个</Form.Text>
                </Form.Item>
                <Form.Item label="最大节点数">
                  <Form.Text>{licenseConfig.cluster.maxNodeNum} 个</Form.Text>
                </Form.Item>
                <Form.Item label="最大核数">
                  <Form.Text>{licenseConfig.cluster.maxCpuNum} 个</Form.Text>
                </Form.Item>
                <Form.Item label="是否支持GPU虚拟化">
                  <Form.Text>
                    {licenseConfig.cluster.gpuShare ? "是" : "否"}
                  </Form.Text>
                </Form.Item>
              </Form>
              <hr />
              <Form.Title>Ceph限制</Form.Title>
              <Form layout="fixed">
                <Form.Item label="最大节点数">
                  <Form.Text>{licenseConfig.ceph.maxNodeNum} 个</Form.Text>
                </Form.Item>
                <Form.Item label="最大存储空间">
                  <Form.Text>{licenseConfig.ceph.maxStorage} GB</Form.Text>
                </Form.Item>
              </Form>
              <hr />

              <Form.Title>ElasticSearch限制</Form.Title>
              <Form layout="fixed">
                <Form.Item label="最大节点数">
                  <Form.Text>{licenseConfig.ES.maxNodeNum} 个</Form.Text>
                </Form.Item>
                <Form.Item label="最大存储空间">
                  <Form.Text>{licenseConfig.ES.maxStorage} GB</Form.Text>
                </Form.Item>
              </Form>
              <Form.Action>
                <Button
                  type="primary"
                  className="mr10"
                  onClick={() => {
                    actions.installer.stepNext(1);
                  }}
                >
                  上一步
                </Button>
                <Button
                  type="primary"
                  onClick={() => {
                    actions.installer.stepNext();
                  }}
                >
                  下一步
                </Button>
              </Form.Action>
            </div>
          </Card.Body>
        </Card>
      </div>
    ) : (
      <noscript />
    );
  }
}
