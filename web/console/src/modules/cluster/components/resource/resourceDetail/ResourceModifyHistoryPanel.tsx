import * as yaml from 'js-yaml';
import * as React from 'react';
import { connect } from 'react-redux';

import { Bubble, TableColumn, Text } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { dateFormatter } from '../../../../../../helpers';
import { GridTable, LinkButton, TipDialog } from '../../../../common/components';
import { allActions } from '../../../actions';
import { Replicaset } from '../../../models';
import { RootProps } from '../../ClusterApp';
import { YamlEditorPanel } from '../YamlEditorPanel';

interface ResourceModifyHistoryPanelState {
  /** 是否展示yaml弹窗 */
  isShowYamlContent?: boolean;

  /** 当前操作的列表的id */
  rsId?: string;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourceModifyHistoryPanel extends React.Component<RootProps, ResourceModifyHistoryPanelState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      isShowYamlContent: false,
      rsId: ''
    };
  }

  componentDidMount() {
    let { actions, subRoot, route } = this.props,
      { resourceSelection } = subRoot.resourceOption;

    // 这里主要是几个tab之间相互跳，需要重新拉取
    resourceSelection.length > 0 &&
      actions.resourceDetail.rs.applyFilter({
        regionId: +route.queries['rid'],
        clusterId: route.queries['clusterId'],
        namespace: route.queries['np']
      });
  }

  componentWillReceiveProps(nextProps: RootProps) {
    let resourceSelection = this.props.subRoot.resourceOption.resourceSelection,
      nextResourceSelection = nextProps.subRoot.resourceOption.resourceSelection;

    if (resourceSelection.length === 0 && nextResourceSelection.length) {
      let { route, actions } = nextProps;
      // 拉取rs的列表
      actions.resourceDetail.rs.applyFilter({
        regionId: +route.queries['rid'],
        clusterId: route.queries['clusterId'],
        namespace: route.queries['np']
      });
    }
  }

  render() {
    return (
      <React.Fragment>
        {this._renderTablePanel()}
        {this._renderContentYamlDialog()}
      </React.Fragment>
    );
  }

  /** 列表的展示 */
  private _renderTablePanel() {
    let { subRoot, actions, route } = this.props,
      { resourceDetailState } = subRoot,
      { rsList, rsQuery } = resourceDetailState;

    const columns: TableColumn<Replicaset>[] = [
      {
        key: 'version',
        header: t('版本号'),
        width: '10%',
        render: x => (
          <Text parent="div" overflow>
            {'v' + x.metadata.annotations['deployment.kubernetes.io/revision']}
          </Text>
        )
      },
      {
        key: 'content',
        header: t('版本详情'),
        width: '10%',
        render: x => (
          <div>
            <i
              className="icon-log"
              style={{ cursor: 'pointer' }}
              data-title={t('查看YAML')}
              data-logviewer
              onClick={() => {
                this.setState({ isShowYamlContent: true, rsId: x.id + '' });
              }}
            />
          </div>
        )
      },
      {
        key: 'registry',
        header: t('镜像'),
        width: '25%',
        render: x => this._renderImageInfo(x)
      },
      {
        key: 'updateTime',
        header: t('更新时间'),
        width: '15%',
        render: x => (
          <Text parent="div" overflow>
            {dateFormatter(new Date(x.metadata.creationTimestamp), 'YYYY-MM-DD HH:mm:ss')}
          </Text>
        )
      },
      {
        key: 'operator',
        header: t('操作'),
        width: '10%',
        render: x => this._renderOperationButton(x)
      }
    ];

    return (
      <GridTable
        columns={columns}
        emptyTips={<div>{t('修订历史列表为空')}</div>}
        listModel={{
          list: rsList,
          query: rsQuery
        }}
        actionOptions={actions.resourceDetail.rs}
      />
    );
  }

  /** 展示镜像的相关信息 */
  private _renderImageInfo(rs: Replicaset) {
    let { containers } = rs.spec.template.spec;

    /**
     * 镜像的形式
     * nginx:latest
     * domain:port/ns/repo:tag
     */
    const getRegistryInfo = (image: string) => {
      let imageInfo: any[] = image.split('/'),
        registry = '',
        tag = '';

      if (imageInfo.length === 1) {
        let [rName, tagName] = imageInfo[0].split(':');
        registry = rName;
        tag = tagName ? tagName : '';
      } else {
        let tagInfo = imageInfo[imageInfo.length - 1].split(':');
        registry = imageInfo.slice(0, imageInfo.length - 1).join('/') + '/' + tagInfo[0];
        tag = tagInfo[1] ? tagInfo[1] : '';
      }
      return [registry, tag];
    };

    let [registry, tag] = getRegistryInfo(containers[0] ? containers[0].image : '');

    return (
      <Bubble
        placement="top-start"
        content={
          containers.map((container, index) => {
            let [registry, tag] = getRegistryInfo(container ? container.image : '');
            return (
              <p key={index} className="text-overflow" style={{ fontSize: '12px' }}>
                <span style={{ display: 'block' }}>{`容器名称：${container.name}`}</span>
                <span style={{ display: 'block' }}>{`镜像：${registry}`}</span>
                <span style={{ display: 'block' }}>{`版本（tag）：${tag}`}</span>
                {containers.length > 1 && <br />}
              </p>
            );
          }) || null
        }
      >
        <p className="text-overflow" style={{ fontSize: '12px' }}>
          <span style={{ display: 'block' }}>{`镜像：${registry}`}</span>
          <span style={{ display: 'block' }}>{`版本（tag）：${tag}`}</span>
        </p>
        {containers.length > 1 && <p className="text">{`等${containers.length}多个镜像`}</p>}
      </Bubble>
    );
  }

  /** 生成操作按钮 */
  private _renderOperationButton(rs: Replicaset) {
    let { resourceOption, resourceDetailState } = this.props.subRoot,
      { rsList } = resourceDetailState,
      { resourceSelection } = resourceOption;

    let disabled = false,
      errorTip = '';

    if (rsList.data.recordCount) {
      let rsIndex = rsList.data.records.findIndex(item => item.id === rs.id);
      if (rsIndex === 0) {
        disabled = true;
        errorTip = t('不能回滚至当前版本');
      }
    }

    return (
      <LinkButton
        tipDirection={'right'}
        errorTip={errorTip}
        disabled={disabled}
        onClick={() => {
          if (!disabled) {
            this._handleRollbackOperation(rs);
          }
        }}
      >
        {t('回滚')}
      </LinkButton>
    );
  }

  /** 处理回滚操作 */
  private _handleRollbackOperation(rs: Replicaset) {
    let { actions } = this.props;

    // 选择rs
    actions.resourceDetail.rs.selectRs([rs]);
    // 弹出确认框
    actions.workflow.rollbackResource.start([]);
  }

  /** 展示版本的具体yaml内容 */
  private _renderContentYamlDialog() {
    let { isShowYamlContent, rsId } = this.state,
      { rsList } = this.props.subRoot.resourceDetailState;

    let rsInfo = rsList.data.records.find(item => item.id === rsId);

    let content = rsInfo ? yaml.safeDump(rsInfo) : '';
    return (
      <TipDialog
        caption={t('版本详情')}
        width={700}
        isShow={isShowYamlContent}
        performAction={() => this.setState({ isShowYamlContent: false, rsId: '' })}
        cancelAction={() => this.setState({ isShowYamlContent: false, rsId: '' })}
      >
        <YamlEditorPanel config={content} readOnly={true} height={450} />
      </TipDialog>
    );
  }
}
