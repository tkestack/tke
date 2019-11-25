import * as React from "react";
import { RootProps } from "./InstallerApp";
import { Button, Bubble, Card, Input, Text } from "@tencent/tea-component";

export class WelcomePanel extends React.Component<RootProps> {
  render() {
    const { isVerified, actions, editState } = this.props;
    return (
      <div style={{ maxWidth: "1000px", minHeight: "600px", margin: "0 auto" }}>
        <h2 style={{ margin: "40px 0px", fontWeight: 600 }}>TKE Enterprise</h2>
        <Card>
          <Card.Body>
            <div
              style={{
                padding: "60px 60px 20px",
                fontSize: "14px",
                backgroundColor: "#fff"
              }}
            >
              <h2>欢迎使用腾讯云容器服务TKE企业版</h2>
              <p style={{ margin: "20px 0", fontSize: "14px" }}>
                请阅读以下安装事宜后开始安装部署腾讯云容器服务TKE:
              </p>
              <p style={{ lineHeight: "1.8" }}>
                1. 请提前准备能与本机连通的TKE企业版的支撑设备。
              </p>
              <p style={{ lineHeight: "1.8" }}>
                2. 安装完成后将返回TKE系统API
                Server根证书和超级管理员账号密码，请妥善保存。
              </p>
              <p style={{ lineHeight: "1.8" }}>
                3. 输入License后点击下方开始按钮，进入安装部署。
              </p>
              <p style={{ marginTop: "40px", lineHeight: "1.8" }}>
                请输入TKE Enterprise License:
              </p>
              <div>
                <Input
                  style={{ width: "840px", height: "120px" }}
                  multiline
                  value={editState.license}
                  onChange={value =>
                    actions.installer.updateEdit({ license: value })
                  }
                />
                <Text
                  theme={
                    isVerified === 0
                      ? "danger"
                      : isVerified === 1
                      ? "success"
                      : "text"
                  }
                >
                  {isVerified === 0
                    ? "License 无效，请输入正确的License"
                    : isVerified === 1
                    ? "License 验证通过"
                    : ""}
                </Text>
              </div>
              <div style={{ marginTop: "20px" }}>
                <Bubble
                  content={!editState.license ? "请先输入正确的License" : ""}
                >
                  <Button
                    disabled={!editState.license}
                    type="primary"
                    onClick={() => {
                      actions.installer.verifyLicense(editState.license);
                    }}
                  >
                    开始
                  </Button>
                </Bubble>
              </div>
            </div>
          </Card.Body>
        </Card>
      </div>
    );
  }
}
