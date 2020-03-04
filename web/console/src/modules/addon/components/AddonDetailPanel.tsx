import * as React from 'react';
import { RootProps } from './AddonApp';
import { router } from '../router';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { connect } from 'react-redux';
import { allActions } from '../actions';
import { Resource } from '../../common';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Text, Icon } from '@tencent/tea-component';
import { AddonStatusNameMap, AddonStatusThemeMap, AddonTypeMap } from '../constants/Config';
import { dateFormatter } from '../../../../helpers';
import { FetchState } from '@tencent/ff-redux';
import { FormPanel } from '@tencent/ff-component';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class AddonDetailPanel extends React.Component<RootProps, {}> {
  render() {
    return this._renderBasicInfo();
  }

  /** 展示基础数据 */
  private _renderBasicInfo() {
    let { openAddon } = this.props;

    let content: React.ReactNode;

    if (
      openAddon.list.fetched !== true ||
      openAddon.list.fetchState === FetchState.Fetching ||
      openAddon.selection === null
    ) {
      content = <Icon type="loading" />;
    } else {
      let addonInfo: Resource = openAddon.selection;

      let status = addonInfo.status.phase.toLowerCase() || '-';
      let theme = AddonStatusThemeMap[status];

      // 创建时间
      let time: any = '-';
      if (addonInfo.metadata.creationTimestamp) {
        time = dateFormatter(new Date(addonInfo.metadata.creationTimestamp), 'YYYY-MM-DD HH:mm:ss');
      }

      content = (
        <React.Fragment>
          <FormPanel.Item text label={t('组件名称')}>
            <Text>{addonInfo.metadata.name || '-'}</Text>
          </FormPanel.Item>
          <FormPanel.Item text label={t('来源')}>
            <Text>{addonInfo.spec.type}</Text>
          </FormPanel.Item>
          <FormPanel.Item text label={t('状态')}>
            <Text theme={theme}>{AddonStatusNameMap[status]}</Text>
          </FormPanel.Item>
          <FormPanel.Item text label={t('类型')}>
            <Text>{AddonTypeMap[addonInfo.spec.level || 'Basic']}</Text>
          </FormPanel.Item>
          <FormPanel.Item text label={t('版本')}>
            <Text>{addonInfo.spec.version || '-'}</Text>
          </FormPanel.Item>
          <FormPanel.Item text label={t('创建时间')}>
            <Text>{time}</Text>
          </FormPanel.Item>
        </React.Fragment>
      );
    }

    return <FormPanel title={t('基本信息')}>{content}</FormPanel>;
  }
}
