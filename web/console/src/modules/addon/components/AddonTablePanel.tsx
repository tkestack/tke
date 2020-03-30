import * as React from 'react';
import { connect } from 'react-redux';

import { TablePanel, TablePanelColumnProps } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Bubble, Icon, Text } from '@tencent/tea-component';

import { dateFormatter } from '../../../../helpers';
import { Clip, LinkButton, Resource } from '../../common';
import { allActions } from '../actions';
import {
  AddonStatusEnum,
  AddonStatusNameMap,
  AddonStatusThemeMap,
  AddonTypeEnum,
  AddonTypeMap
} from '../constants/Config';
import { RootProps } from './AddonApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class AddonTablePanel extends React.Component<RootProps, {}> {
  render() {
    return this._renderTablePanel();
  }

  /** 展示列表 */
  private _renderTablePanel() {
    let { openAddon, actions, addon, route } = this.props;

    const columns: TablePanelColumnProps<Resource>[] = [
      {
        key: 'name',
        header: t('名称/来源'),
        width: '15%',
        render: x => (
          <React.Fragment>
            <Text id={x.metadata.name} overflow verticalAlign="middle">
              {x.metadata.name || '-'}
            </Text>
            <Clip target={`#${x.metadata.name}`} />
            <Text parent="p">{x.spec.type || '-'}</Text>
          </React.Fragment>
        )
      },
      {
        key: 'status',
        header: t('状态'),
        width: '10%',
        render: x => {
          // 这里进行小写格式化是因为后台返回的status 大小写不一定统一
          let status = x.status.phase.toLowerCase() || '-';
          let theme = AddonStatusThemeMap[status] || 'text';
          // let finder = addon.list.data.records.find(item => item.type === x.spec.type);
          let content: React.ReactNode = <noscript />;

          let basicContent: React.ReactNode = (
            <Text theme={theme} verticalAlign="middle">
              {AddonStatusNameMap[status]}
            </Text>
          );

          // 先展示错误，再展示可更新，错误影响功能
          if (status === AddonStatusEnum.Failed) {
            content = (
              <React.Fragment>
                {basicContent}
                <Bubble content={<Text>{x.status.reason || '-'}</Text>}>
                  <Icon type="error" />
                </Bubble>
              </React.Fragment>
            );
          } else if (status !== AddonStatusEnum.Running) {
            content = (
              <React.Fragment>
                {basicContent}
                <Icon type="loading" />
                {x.status.reason && (
                  <Bubble content={<Text>{x.status.reason || '-'}</Text>}>
                    <Icon type="error" />
                  </Bubble>
                )}
              </React.Fragment>
            );
            // } else if (status === AddonStatusEnum.Running && finder && finder.latestVersion !== x.spec.version) {
            // 注意判断条件里面，版本的判断，需要去除v，然后进行比较
            //   theme = 'warning';
            //   content = (
            //     <React.Fragment>
            //       {basicContent}
            //       <Bubble content={<Text>{`可升级，最新版本为 ${finder.latestVersion}`}</Text>}>
            //         <Icon type="warning" />
            //       </Bubble>
            //     </React.Fragment>
            //   );
          } else {
            content = basicContent;
          }
          return content;
        }
      },
      {
        key: 'type',
        header: t('类型'),
        width: '15%',
        render: x => <Text overflow>{AddonTypeMap[x.spec.level || AddonTypeEnum.Basic]}</Text>
      },
      {
        key: 'version',
        header: t('版本'),
        width: '12%',
        render: x => <Text>{x.spec.version || '-'}</Text>
      },
      {
        key: 'createTime',
        header: t('创建时间'),
        width: '15%',
        render: x => {
          let time: any = '-';
          if (x.metadata.creationTimestamp) {
            time = dateFormatter(new Date(x.metadata.creationTimestamp), 'YYYY-MM-DD HH:mm:ss');
          }
          return <Text>{time}</Text>;
        }
      }
    ];

    return (
      <TablePanel
        columns={columns}
        action={actions.cluster.addon}
        model={openAddon}
        emptyTips={t('当前集群下尚未安装扩展组件')}
        getOperations={x => this._renderOperationButtons(x)}
        operationsWidth={180}
      />
    );
  }

  /** 构建操作列的展示 */
  private _renderOperationButtons(clusterAddon: Resource) {
    let { addon, actions } = this.props;
    let isCanNotDelete = clusterAddon.spec.level === AddonTypeEnum.Basic;

    let finder = addon.list.data.records.find(item => item.type === clusterAddon.spec.type);
    let isCanUpdateAddon = finder ? finder.latestVersion !== clusterAddon.spec.version : false;

    const renderDeleteButton = () => {
      return (
        <LinkButton
          key={0}
          // disabled={isCanNotDelete}
          errorTip={t('基础组件不可删除')}
          tipDirection="left"
          onClick={() => {
            // 选择当前项
            actions.cluster.addon.select(clusterAddon);
            actions.workflow.deleteResource.start();
          }}
        >
          {t('删除')}
        </LinkButton>
      );
    };

    let buttons: any[] = [];
    // buttons.push(renderUpdateButton());
    buttons.push(renderDeleteButton());

    return buttons;
  }
}
