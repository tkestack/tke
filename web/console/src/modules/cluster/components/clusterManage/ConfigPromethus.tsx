import React from 'react';
import { FormPanel } from '@tencent/ff-component';
import { Button, ContentView, Justify, Icon } from '@tencent/tea-component';
import { getWorkflowError, InputField, ResourceInfo, TipInfo } from '../../../../modules/common';
import { t } from '@tencent/tea-app/lib/i18n';

export function ConfigPromethus() {
  return (
    <ContentView>
      <ContentView.Header>
        <Justify
          left={
            <React.Fragment>
              <a href="javascript:;" className="back-link" onClick={history.back}>
                <Icon type="btnback" />
                {t('返回')}
              </a>
              <span className="line-icon">|</span>
              <h2>集群监控配置</h2>
            </React.Fragment>
          }
        />
      </ContentView.Header>
      <ContentView.Body>
        <FormPanel>
          <FormPanel.Item label={t('名称')}>
            <InputField type="text" value={name} placeholder={t('请输入集群名称')} tipMode="popup" />
          </FormPanel.Item>

          <FormPanel.Footer>
            <React.Fragment>
              <Button className="m" type="primary">
                {t('提交')}
              </Button>
              <Button type="weak">取消</Button>
              <TipInfo type="error" isForm>
                {/* {getWorkflowError(workflow)} */}
              </TipInfo>
            </React.Fragment>
          </FormPanel.Footer>
        </FormPanel>
      </ContentView.Body>
    </ContentView>
  );
}
