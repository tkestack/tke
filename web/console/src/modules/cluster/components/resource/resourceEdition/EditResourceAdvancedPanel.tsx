import * as React from 'react';
import { connect } from 'react-redux';

import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Select } from '@tencent/tea-component';

import { FormItem } from '../../../../common';
import { allActions } from '../../../actions';
import {
    FloatingIPReleasePolicy, WorkloadNetworkType, WorkloadNetworkTypeEnum
} from '../../../constants/Config';
import { RootProps } from '../../ClusterApp';
import { EditResourceAnnotations } from './EditResourceAnnotations';
import { EditResourceImagePullSecretsPanel } from './EditResourceImagePullSecretsPanel';
import { EditResourceNodeAffinityPanel } from './EditResourceNodeAffinityPanel';

interface EditResourceAdvancedPanelProps extends RootProps {
  /** 是否展示高级设置 */
  isOpenAdvanced: boolean;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditResourceAdvancedPanel extends React.Component<EditResourceAdvancedPanelProps, {}> {
  render() {
    let { isOpenAdvanced, subRoot, actions } = this.props,
      { workloadEdit } = subRoot,
      { networkType, floatingIPReleasePolicy } = workloadEdit;

    // let isShowPort = networkType !== 'Overlay';

    return isOpenAdvanced ? (
      <React.Fragment>
        <EditResourceImagePullSecretsPanel />
        <EditResourceNodeAffinityPanel />
        <EditResourceAnnotations />
        <FormItem label={t('网络模式')}>
          <Select
            size="m"
            options={WorkloadNetworkType}
            value={networkType}
            onChange={value => {
              actions.editWorkload.selectNetworkType(value);
            }}
          />
        </FormItem>
        <FormItem isShow={networkType === WorkloadNetworkTypeEnum.FloatingIP} label={t('IP回收策略')}>
          <FormPanel.Select
            size="m"
            options={FloatingIPReleasePolicy}
            value={floatingIPReleasePolicy}
            onChange={value => {
              actions.editWorkload.selectFloatingIPReleasePolicy(value);
            }}
          ></FormPanel.Select>
        </FormItem>
        {/* <FormItem label={t('端口')} isShow={isShowPort}>
        </FormItem> */}
      </React.Fragment>
    ) : (
      <noscript />
    );
  }
}
