import * as React from 'react';
import { bindActionCreators, uuid } from '@tencent/qcloud-lib';
import { Bubble, Table, TableColumn, Text } from '@tea/component';
import { connect } from 'react-redux';
import { allActions } from '../../../actions';
import { RootProps } from '../../ClusterApp';
import { DetailLayout } from '../../../../common/layouts';
import { ListItem, LinkButton } from '../../../../common/components';
import * as classnames from 'classnames';
import { ResourceStatus } from '../../../constants/Config';
import { isEmpty } from '../../../../common/utils';
import { dateFormatter } from '../../../../../../helpers';
import { FetchState } from '@tencent/qcloud-redux-fetcher';
import { Resource } from '../../../models';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { stylize } from '@tea/component/table/addons/stylize';

const loadingElement: JSX.Element = (
  <div>
    <i className="n-loading-icon" />
    &nbsp; <span className="text">{t('加载中...')}</span>
  </div>
);

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
export class ResourceNamespaceDetailPanel extends React.Component<RootProps, {}> {
  render() {
    let { subRoot } = this.props,
      { resourceOption } = subRoot,
      { resourceSelection, resourceList } = resourceOption;

    let resourceIns = resourceSelection[0];

    let statusMap = ResourceStatus['np'];

    let isNeedLoading =
      resourceList.fetched !== true ||
      resourceList.fetchState === FetchState.Fetching ||
      resourceSelection.length === 0;

    return isNeedLoading ? (
      loadingElement
    ) : (
      <div>
        <DetailLayout>
          <div className="param-box">
            <div className="param-hd">
              <h3>{t('基本信息')}</h3>
            </div>
            <div className="param-bd">
              <ul className="item-descr-list">
                <ListItem label={t('名称')}>{resourceIns.metadata.name}</ListItem>

                <ListItem label={t('状态')}>
                  <p
                    className={classnames(
                      '',
                      statusMap[resourceIns.status.phase] && statusMap[resourceIns.status.phase].classname
                    )}
                  >
                    {(statusMap[resourceIns.status.phase] && statusMap[resourceIns.status.phase].text) || '-'}
                  </p>
                </ListItem>

                <ListItem label={t('描述')}>
                  {resourceIns.metadata.annotations ? resourceIns.metadata.annotations.description : '-'}
                </ListItem>

                <ListItem label={t('创建时间')}>
                  {dateFormatter(new Date(resourceIns.metadata.creationTimestamp), 'YYYY-MM-DD HH:mm:ss')}
                </ListItem>
              </ul>
            </div>
          </div>
        </DetailLayout>
      </div>
    );
  }
}
