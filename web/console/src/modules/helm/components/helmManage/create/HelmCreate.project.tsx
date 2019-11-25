import * as React from 'react';
import { RootProps } from '../../HelmApp';
import { router } from '../../../router';
import { MainBodyLayout, FormLayout } from '../../../../common/layouts';
import { TipInfo } from '../../../../common/components';
import { Button } from '@tea/component/button';
import { FetchState } from '@tencent/qcloud-redux-fetcher';
import { HelmResource, OtherType } from '../../../constants/Config';
import { BaseInfoPanel } from './BaseInfoPanel';
import { TencentHubChartPanel } from './TencentHubChartPanel';
import { OtherChartPanel } from './OtherChartPanel';
import { KeyValuesPanel } from './KeyValuesPanel';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

export class HelmCreate extends React.Component<RootProps, {}> {
  componentDidMount() {
    // this.props.actions.create.clear();
    this.props.actions.create.inputKeyValue([]);
  }
  componentWillUnmount() {
    this.props.actions.create.clear();
    this.props.actions.create.inputKeyValue([]);
    // 去除错误信息
  }
  goBack() {
    let { actions, route } = this.props,
      urlParams = router.resolve(route);
    router.navigate({}, route.queries);
  }
  onOk() {
    const {
      actions,
      helmCreation: { name }
    } = this.props;
    actions.create.validAll();
    let canSave = this.isCanSave();
    if (canSave) {
      this.props.actions.create.createHelm();
    }
  }
  onCancel() {
    this.goBack();
  }

  isCanSave() {
    let {
      namespaceSelection,
      helmCreation: { isValid, resourceSelection, otherTypeSelection }
    } = this.props;
    let canSave = true;
    if (isValid.name !== '' || namespaceSelection === '') {
      canSave = false;
    } else {
      if (resourceSelection === HelmResource.Other) {
        if (isValid.otherChartUrl !== '') {
          canSave = false;
        } else {
          if (otherTypeSelection === OtherType.Private) {
            if (isValid.otherUserName !== '' || isValid.otherPassword !== '') {
              canSave = false;
            }
          }
        }
      }
    }
    return canSave;
  }
  render() {
    let {
      helmCreation: { resourceSelection }
    } = this.props;

    let canSave = this.isCanSave();
    return (
      <div>
        <div className="manage-area-title secondary-title">
          <a
            href="javascript:void(0)"
            className="back-link"
            onClick={() => {
              this.goBack();
            }}
          >
            <i className="btn-back-icon" />
            <span>{t('返回')}</span>
          </a>
          <span className="line-icon"> |</span>
          <h2>{t('新建 Helm 应用')}</h2>
        </div>
        <MainBodyLayout>
          <div className="manage-area-main-inner">
            {/* <TipInfo>
              <span style={{ verticalAlign: 'middle' }}>
                {t(
                  '创建Helm应用，若应用中包含了公网CLB类型的Services或Ingress，将按照腾讯云CLB对应价格收费。若应用中包含PV/PVC/StorageClass，其创建的存储资源将按对应的产品价格收费。'
                )}
              </span>
            </TipInfo> */}
            <FormLayout>
              <div className="param-box">
                <div className="param-bd">
                  <BaseInfoPanel {...this.props} />
                  {resourceSelection === HelmResource.TencentHub && <TencentHubChartPanel {...this.props} />}
                  {resourceSelection === HelmResource.Other && <OtherChartPanel {...this.props} />}
                  <KeyValuesPanel
                    onChangeKeyValue={this.props.actions.create.inputKeyValue}
                    kvs={this.props.helmCreation.kvs}
                  />
                </div>
                <div className="param-ft">
                  <Button
                    className="mr10"
                    title={t('完成')}
                    disabled={!canSave}
                    onClick={() => {
                      this.onOk();
                    }}
                    type="primary"
                  >
                    {t('完成')}
                  </Button>
                  <Button
                    title={t('取消')}
                    onClick={() => {
                      this.onCancel();
                    }}
                  >
                    {t('取消')}
                  </Button>
                </div>
              </div>
            </FormLayout>
          </div>
        </MainBodyLayout>
      </div>
    );
  }
}
