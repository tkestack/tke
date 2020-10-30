import classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators, FetchState } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Icon, Text, Form } from '@tencent/tea-component';

import { InputField, LinkButton, SelectList } from '../../common/components';
import { cloneDeep } from '../../common/utils/cloneDeep';
import { allActions } from '../actions';
import { logModeList, ResourceListMapForPodLog } from '../constants/Config';
import { router } from '../router';
import { RootProps } from './LogStashApp';

export interface PropTypes extends RootProps {
  isEdit?: boolean; // 是否是编辑模式
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(
  state => state,
  mapDispatchToProps
)
export class EditOriginContainerFilePanel extends React.Component<PropTypes, {}> {
  render() {
    let { actions, logStashEdit, route, clusterList, logSelection, namespaceList } = this.props,
      {
        logMode,
        v_containerFileNamespace,
        v_containerFileWorkload,
        v_containerFileWorkloadType,
        containerFileNamespace,
        containerFileWorkload,
        containerFileWorkloadType,
        containerFileWorkloadList,
        resourceList,
        podList
      } = logStashEdit;
    let mode = router.resolve(route)['mode'];
    let isNamespaceNeedLoading =
      clusterList.fetchState === FetchState.Fetching || namespaceList.fetchState === FetchState.Fetching;
    let isWorkloadNeedLoading =
      clusterList.fetchState === FetchState.Fetching ||
      resourceList.fetchState === FetchState.Fetching ||
      isNamespaceNeedLoading;

    //是否可以更改容器文件路径的输出
    let ifCanUpdateContainerFile = logSelection[0] && logSelection[0].inputType === 'pod-log' && mode === 'update';
    return (
      <FormPanel.Item label={t('日志源')} isShow={logMode === logModeList.containerFile.value}>
        <FormPanel isNeedCard={false} fixed style={{ minWidth: 600, padding: '30px ' }}>
          <FormPanel.Item label={t('工作负载选项')}>
            <div
              className={classnames('code-list', {
                'is-error': v_containerFileNamespace.status === 2
              })}
              style={{ display: 'inline-block' }}
            >
              {this.props.isEdit ? (
                <Form.Text
                  style={{
                    // width: '10px',
                    marginRight: '5px'
                  }}
                >
                  {containerFileNamespace}
                </Form.Text>
              ) : (
                <SelectList
                  recordData={namespaceList}
                  className="tc-15-select "
                  valueField="namespaceValue"
                  textField="namespace"
                  onSelect={value => {
                    actions.editLogStash.selectContainerFileNamespace(value);
                    actions.namespace.selectNamespace(value);
                    // 兼容业务侧的处理
                    if (window.location.href.includes('tkestack-project')) {
                      let namespaceFound = namespaceList.data.records.find(item => item.namespaceValue === value);
                      actions.cluster.selectClusterFromEditNamespace(namespaceFound.cluster);
                    }
                  }}
                  name="Namespace"
                  value={containerFileNamespace}
                  disabled={ifCanUpdateContainerFile || namespaceList.data.recordCount === 0}
                  style={{
                    width: '180px',
                    marginRight: '5px'
                  }}
                />
              )}
            </div>

            <div
              className={classnames('code-list', {
                'is-error': v_containerFileWorkloadType.status === 2
              })}
              style={{ display: 'inline-block' }}
            >
              <SelectList
                name="工作负载类型"
                valueField="value"
                textField="name"
                recordList={ResourceListMapForPodLog}
                className="tc-15-select "
                onSelect={value => {
                  actions.editLogStash.selectContainerFileWorkloadType(value);
                }}
                value={containerFileWorkloadType}
                disabled={ifCanUpdateContainerFile}
                style={{
                  width: '180px',
                  marginRight: '5px'
                }}
              />
            </div>
            {isWorkloadNeedLoading ? (
              <Icon type="loading" style={{ marginRight: '174px' }} />
            ) : (
              <div
                className={classnames('code-list', {
                  'is-error': v_containerFileWorkload.status === 2
                })}
                style={{ display: 'inline-block' }}
              >
                <SelectList
                  name="workload"
                  recordList={containerFileWorkloadList}
                  valueField="name"
                  textField="name"
                  onSelect={value => {
                    actions.editLogStash.selectContainerFileWorkload(value);
                  }}
                  disabled={containerFileWorkloadList.length === 0 || ifCanUpdateContainerFile}
                  className="tc-15-select "
                  value={containerFileWorkload}
                  style={{
                    width: '180px',
                    marginRight: '5px'
                  }}
                />
              </div>
            )}
          </FormPanel.Item>
          {containerFileWorkload && containerFileNamespace && containerFileWorkloadType ? (
            <FormPanel.Item
              label={t('配置采集路径')}
              message={
                <React.Fragment>
                  <Text parent="p">
                    <Trans>文件路径若输入`stdout`，则转为容器标准输出模式</Trans>
                  </Text>
                  <Text parent="p">
                    <Trans>
                      可配置多个路径。路径必须以`/`开头和结尾，文件名支持通配符（*）。文件路径和文件名最长支持63个字符
                    </Trans>
                  </Text>
                  <Text parent="p">
                    <Trans>请保证容器的日志文件保存在数据卷，否则收集规则无法生效 </Trans>
                  </Text>
                </React.Fragment>
              }
            >
              {this._renderPathDataList()}
              <LinkButton onClick={actions.editLogStash.addContainerFileContainerFilePath}>添加路径</LinkButton>
            </FormPanel.Item>
          ) : (
            <noscript />
          )}
        </FormPanel>
      </FormPanel.Item>
    );
  }

  private _renderPathDataList() {
    const { actions, logStashEdit } = this.props;
    const { containerFilePaths, podList } = logStashEdit;

    //因为原始的podList的需要显示的字段在spec字段的container字段里的name字段里，并没有直接在直接字段里面，所以直接拿来podList的recodeData展示展示不出来，需要做修改
    let selectPodList = cloneDeep(podList);
    selectPodList.data.records = [];
    podList.data.records.forEach(pod => {
      pod.spec.containers.forEach(container => {
        selectPodList.data.records.push({
          name: container.name
        });
      });
    });

    return containerFilePaths.map((path, index) => {
      return (
        <React.Fragment key={index}>
          <div style={{ display: 'flex', alignItems: 'center', marginBottom: '5px' }}>
            <div style={{ paddingRight: '30px' }}>
              <label className="text-label" style={{ marginRight: '15px', verticalAlign: '-5px' }}>
                {t('容器名')}
              </label>
              <div
                className={classnames('code-list', {
                  'is-error': path.v_containerName.status === 2
                })}
                style={{ display: 'inline-block' }}
              >
                <SelectList
                  className="tc-15-select m"
                  recordData={selectPodList}
                  name="容器名"
                  style={{ display: 'inline-block', width: '190px' }}
                  onSelect={value => {
                    actions.editLogStash.selectContainerFileContainerName(value, index);
                  }}
                  value={path.containerName}
                  valueField="name"
                  textField="name"
                  disabled={podList.data.recordCount === 0}
                />
              </div>
            </div>
            <div style={{ paddingRight: '20px' }}>
              <label className="text-label" style={{ marginRight: '15px' }}>
                {t('文件路径')}
              </label>
              <InputField
                type="text"
                placeholder={t('如/data/log/logfilename-*.log')}
                value={path.containerFilePath}
                tipMode="popup"
                validator={path.v_containerFilePath}
                onChange={value => {
                  actions.editLogStash.inputContainerFileContainerFilePath(value, index);
                }}
                style={{ width: '190px' }}
                onBlur={() => {
                  actions.validate.validateContainerFileContainerFilePath(index);
                }}
              />
            </div>
            <div>
              <LinkButton
                onClick={() => actions.editLogStash.deleteContainerFileContainerFilePath(index)}
                disabled={containerFilePaths.length <= 1}
                errorTip={t('至少要有一个文件路径')}
                tip={t('删除')}
                tipDirection="right"
                isShow={index > 0}
                className="inline-help-text"
              >
                <Icon type="close" />
              </LinkButton>
            </div>
          </div>
        </React.Fragment>
      );
    });
  }
}
