import classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators, insertCSS } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Bubble, Checkbox, FormItem, Form } from '@tencent/tea-component';

import { LinkButton, SelectList } from '../../common/components';
import { cloneDeep } from '../../common/utils';
import { allActions } from '../actions';
import { ResourceListMapForContainerLog } from '../constants/Config';
import { ContainerLogs, Resource } from '../models';
import { isCanAddContainerLog } from './EditOriginContainerPanel';
import { ContainerItemProps } from './ListOriginContainerItemPanel';
import { clusterActions } from '@src/modules/logStash/actions/clusterActions';

insertCSS(
  'EditOriginContainerItemPanel',
  `
  .rich-content .form-ctrl-label {
    margin-top: 0;
  }

  .rich-content .mb20 span {
    display: inline-block;
    vertical-align: middle;
  }

  .rich-content .form-ctrl-label span {
    vertical-align: 1px;
  }
`
);

let loadingElement: JSX.Element = (
  <div>
    <i className="n-loading-icon" />
    &nbsp;{' '}
    <span className="text" style={{ fontSize: '12px' }}>
      {t(' 加载中...')}
    </span>
  </div>
);

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(
  state => state,
  mapDispatchToProps
)
export class EditOriginContainerItemPanel extends React.Component<ContainerItemProps, any> {
  render() {
    let { actions, logStashEdit, cKey, namespaceList } = this.props,
      { containerLogs } = logStashEdit;
    /**
     * 这里传入containerLogIndex 的原因是因为，传入id之后，还是要对 containerLogs进行遍历
     * 这样每次操作，如更改namespace，选择服务等等操作，都需要重复循环数组，效率低
     */
    let containerLogIndex = containerLogs.findIndex(con => con.id === cKey),
      containerLog: ContainerLogs = containerLogs[containerLogIndex];

    let { canSave } = isCanAddContainerLog(containerLogs, namespaceList.data.recordCount),
      canDelete = containerLogs.length > 1;

    //当前可以进行下拉选择的namespace列表，需要剔除已经选择过的
    let optinalNameSpaceList = this._getOptionalNameSpaceList();
    return (
      <FormPanel fixed isNeedCard={false} style={{ minWidth: 600, padding: '30px' }}>
        <div className="run-docker-box" style={containerLog.collectorWay === 'workload' ? { minWidth: '750px' } : {}}>
          <div className="edit-param-list">
            <div className="param-box" style={{ paddingBottom: '0' }}>
              <div className="param-bd">
                <ul className="form-list fixed-layout">
                  <FormPanel.Item label={t('所属Namespace')}>
                    {this.props.isEdit ? (
                      <Form.Text>{containerLog.namespaceSelection}</Form.Text>
                    ) : (
                      <Bubble
                        content={
                          containerLog.v_namespaceSelection.message ? containerLog.v_namespaceSelection.message : null
                        }
                      >
                        <div
                          className={classnames('code-list', {
                            'is-error': containerLog.v_namespaceSelection.status === 2
                          })}
                        >
                          <SelectList
                            value={containerLog.namespaceSelection}
                            recordData={optinalNameSpaceList}
                            valueField="namespaceValue"
                            textField="namespace"
                            className="tc-15-select m"
                            onSelect={value => {
                              actions.editLogStash.selectContainerLogNamespace(value, containerLogIndex);
                              actions.namespace.selectNamespace(value);
                              // 兼容业务侧的处理
                              if (window.location.href.includes('tkestack-project')) {
                                let namespaceFound = namespaceList.data.records.find(item => item.namespaceValue === value);
                                actions.cluster.selectClusterFromEditNamespace(namespaceFound.cluster);
                              }
                            }}
                            name="Namespace"
                            tipPosition="right"
                            style={{
                              display: 'inline-block'
                            }}
                          />
                        </div>
                      </Bubble>
                    )}
                  </FormPanel.Item>

                  {this._renderCollectorWay(containerLog, containerLogIndex)}
                </ul>
              </div>
            </div>
          </div>
        </div>
      </FormPanel>
    );
  }

  /** 处理保存按钮 */
  private _handleSaveContainerLog(canSave: boolean, containerLogIndex: number) {
    let { actions } = this.props;
    // 触发一下校验的动作，检查相关的项是否合法
    actions.validate.validateContainerLog(containerLogIndex);
    if (canSave) {
      actions.editLogStash.updateContainerLog({ status: 'edited' }, containerLogIndex);
    }
  }

  /** 渲染采集对象 */
  private _renderCollectorWay(containerLog: ContainerLogs, containerLogIndex: number) {
    let { actions } = this.props;

    return (
      <FormItem label={t('采集对象')}>
        <div className="form-unit is-success">
          <label className="form-ctrl-label">
            <input
              type="radio"
              style={{ verticalAlign: 'middle' }}
              className="tc-15-radio"
              checked={containerLog.collectorWay === 'container'}
              onChange={e => {
                actions.editLogStash.updateContainerLog({ collectorWay: 'container' }, containerLogIndex);
              }}
            />
            {t('全部容器')}
            <span className="inline-help-text text-label">{t('包含该Namespace下所有的容器')}</span>
          </label>
          <label className="form-ctrl-label">
            <input
              type="radio"
              style={{ verticalAlign: 'middle' }}
              className="tc-15-radio"
              checked={containerLog.collectorWay === 'workload'}
              onChange={e => {
                actions.editLogStash.updateContainerLog({ collectorWay: 'workload' }, containerLogIndex);
              }}
            />
            {t('按工作负载（Workload）选择')}
          </label>

          {containerLog.collectorWay === 'workload' && this._renderCollectorByWorkload(containerLog, containerLogIndex)}
          {containerLog.collectorWay === 'workload' && containerLog.v_workloadSelection.status === 2 && (
            <p className="form-input-help text-danger">{containerLog.v_workloadSelection.message}</p>
          )}
        </div>
      </FormItem>
    );
  }

  /**
   * 获取当前可供选择的namespaceList
   */
  private _getOptionalNameSpaceList() {
    const { logStashEdit, namespaceList } = this.props;
    const { containerLogs } = logStashEdit;

    //获取已经选择了的namespace列表
    let chosenNamespaceList = [];
    for (let i = 0; i < containerLogs.length; i++) {
      if (containerLogs[i].status === 'edited') {
        chosenNamespaceList.push(containerLogs[i].namespaceSelection);
      }
    }
    //筛选还没有选择过的namespaceList
    let result = cloneDeep(namespaceList);
    result.data.records = [];
    for (let i = 0; i < namespaceList.data.recordCount; i++) {
      let canPut = true;
      for (let j = 0; j < chosenNamespaceList.length; j++) {
        if (chosenNamespaceList[j] === namespaceList.data.records[i].namespace) {
          canPut = false;
          break;
        }
      }
      if (canPut) {
        result.data.records.push(namespaceList.data.records[i]);
      }
    }
    return result;
  }

  /** 渲染按照工作负载选择 */
  private _renderCollectorByWorkload(containerLog: ContainerLogs, containerLogIndex: number) {
    let { actions, logStashEdit } = this.props,
      { workloadType, workloadSelection, workloadListFetch, namespaceSelection } = containerLog;

    return (
      <div className="configuration-box" style={{ width: '600px', marginTop: '10px' }}>
        {this._renderLeftSide(containerLog.workloadType, containerLogIndex)}
        <div className="rich-textarea simple-mod" style={{ overflow: 'auto' }}>
          <div className="permission-code-editor">
            <strong className="code-title" style={{ border: 'none' }}>
              {t('列表')}
            </strong>
          </div>
          <div className="rich-content" style={{ padding: '10px 10px 10px 25px' }}>
            {!workloadListFetch[workloadType] && namespaceSelection ? (
              loadingElement
            ) : containerLog['workloadList'][workloadType].length ? (
              <Checkbox.Group
                defaultValue={workloadSelection[workloadType]}
                value={workloadSelection[workloadType]}
                onChange={items => {
                  let workloadSelection = Object.assign({}, containerLog.workloadSelection, {
                    [workloadType]: items
                  });
                  actions.editLogStash.updateContainerLog({ workloadSelection }, containerLogIndex);
                }}
              >
                {containerLog['workloadList'][workloadType].map((workload: Resource) => (
                  <Checkbox
                    name={workload.metadata.name}
                    key={workload.metadata.name}
                    style={{ paddingBottom: '10px' }}
                  >
                    {workload.metadata.name}
                  </Checkbox>
                ))}
                {workloadType}
              </Checkbox.Group>
            ) : (
              <p className="text-label">{t('当前命名空间下，该工作负载类型列表为空')}</p>
            )}
          </div>
        </div>
      </div>
    );
  }

  /** 展示选择框的左半部分 */
  private _renderLeftSide(workloadType: string, containerLogIndex: number) {
    let { actions, route, logStashEdit } = this.props,
      { containerLogs } = logStashEdit;
    let containerLog: ContainerLogs = containerLogs[containerLogIndex];

    return (
      <div className="version-wrap" style={{ width: '150px' }}>
        <div className="tc-15-table-panel version-list">
          <div className="tc-15-table-fixed-head">
            <table className="tc-15-table-box">
              <colgroup>
                <col />
              </colgroup>
              <thead>
                <tr>
                  <th>
                    <div>
                      <span>{t('工作负载类型')}</span>
                    </div>
                  </th>
                </tr>
              </thead>
            </table>
          </div>
          <div className="tc-15-table-fixed-body">
            <table className="tc-15-table-box tc-15-table-rowhover">
              <colgroup>
                <col />
              </colgroup>
              <tbody>
                {ResourceListMapForContainerLog.map((item, index) => {
                  let tdBasicName = item.name;
                  // 还需要展示当前已选的个数/当前工作负载列表的总数,如果用户没有选择namesapce，则不需要展示 未加载
                  tdBasicName +=
                    containerLog.workloadListFetch[item.value] || !containerLog.namespaceSelection
                      ? ` (${containerLog.workloadSelection[item.value].length}/${
                          containerLog.workloadList[item.value].length
                        })`
                      : ` ${t('(未加载)')}`;

                  return (
                    <tr
                      key={index}
                      className={classnames('', { cur: item.value === workloadType })}
                      onClick={() => {
                        let { rid, clusterId } = route.queries;
                        // 更新当前的workload的类型
                        actions.editLogStash.updateContainerLog({ workloadType: item.value }, containerLogIndex);
                        /**
                         * pre: 如果相对应的资源列表已经拉取过了，则切换的时候不需要重复拉取，浪费流量
                         */
                        if (!containerLog.workloadListFetch[item.value]) {
                          actions.resource.applyFilter({
                            workloadType: item.value,
                            clusterId,
                            regionId: +rid,
                            namespace: containerLog.namespaceSelection,
                            isCanFetchResourceList: true
                          });
                        }
                      }}
                    >
                      <td>
                        <div>
                          <span className="text-overflow">{tdBasicName}</span>
                        </div>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    );
  }
}
