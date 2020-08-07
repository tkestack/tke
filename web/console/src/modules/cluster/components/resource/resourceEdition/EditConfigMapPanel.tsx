import * as React from 'react';
import { connect } from 'react-redux';

import { Button, Text } from '@tea/component';
import { bindActionCreators, isSuccessWorkflow, OperationState, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { FormItem, InputField, LinkButton, SelectList, TipInfo } from '../../../../common/components';
import { FormLayout, MainBodyLayout } from '../../../../common/layouts';
import { getWorkflowError } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { validateCMActions } from '../../../actions/validateCMActions';
import { CreateResource } from '../../../models';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';
import { reduceNs } from '../../../../../../helpers';

const ButtonBarStyle = { marginBottom: '5px' };

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditConfigMapPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    let { actions, route } = this.props;
    actions.editCM.selectNamespace(route.queries['np']);
  }

  componentWillUnmount() {
    let { actions } = this.props;
    actions.editCM.clearConfigMapEdit();
  }

  _renderVariableList() {
    let eList = [];
    let { actions, subRoot } = this.props,
      { cmEdit } = subRoot;
    const genVariableTR = (item, index) => {
      return (
        <tr key={index}>
          <td>
            <div className="">
              <InputField
                type="text"
                placeholder={t('请输入名称')}
                className="tc-15-input-text m"
                tipMode="popup"
                validator={item.v_key}
                value={item.key}
                onChange={value => actions.editCM.eidtVariable(item.id, { key: value })}
                onBlur={actions.validate.cm.validateVariableKey}
              />
            </div>
          </td>
          <td>
            <div className="">=</div>
          </td>
          <td>
            <div className="">
              <textarea
                className="tc-15-input-text m"
                style={{ maxWidth: '260px', overflowY: 'visible' }}
                value={item.value}
                onChange={e => actions.editCM.eidtVariable(item.id, { value: e.target.value })}
              />
            </div>
          </td>
          <td>
            <div>
              <a href="javascript:void(0);" onClick={() => actions.editCM.deleteVariable(item.id)}>
                <i className="icon-cancel-icon" />
              </a>
            </div>
          </td>
        </tr>
      );
    };
    cmEdit.variables.forEach((item, index) => {
      eList.push(genVariableTR(item, index));
    });
    return eList;
  }

  render() {
    let { actions, subRoot, route, namespaceList } = this.props,
      urlParams = router.resolve(route),
      { cmEdit, modifyResourceFlow } = subRoot,
      { name, v_name, namespace, variables } = cmEdit;

    /** 加载中的样式 */
    let loadingElement: JSX.Element = (
      <div>
        <i className="n-loading-icon" />
        &nbsp; <span className="text">{t('加载中...')}</span>
      </div>
    );

    /** 渲染namespace列表 */
    let namespaceOptions = namespaceList.data.recordCount
      ? namespaceList.data.records.map((item, index) => {
          return (
            <option key={index} value={item.name}>
              {item.name}
            </option>
          );
        })
      : [];
    namespaceOptions.unshift(
      <option key={-1} value="">
        {namespaceList.data.recordCount ? t('请选择StorageClass') : t('无可用StorageClass')}
      </option>
    );

    let failed = modifyResourceFlow.operationState === OperationState.Done && !isSuccessWorkflow(modifyResourceFlow);

    let canAdd = !cmEdit.variables.some(i => !i.key);

    return (
      <MainBodyLayout>
        <FormLayout>
          <div className="param-box server-update add">
            <ul className="form-list jiqun fixed-layout" style={{ paddingBottom: '50px' }}>
              <FormItem label={t('名称')}>
                <InputField
                  type="text"
                  placeholder={t('请输入名称')}
                  tipMode="popup"
                  validator={v_name}
                  value={name}
                  onChange={actions.editCM.inputCMName}
                  onBlur={actions.validate.cm.validateCMName}
                />
              </FormItem>
              <FormItem label={t('命名空间')}>
                <SelectList
                  className="tc-15-select m"
                  value={namespace}
                  recordData={namespaceList}
                  valueField="name"
                  textField="displayName"
                  name="Namespace"
                  onSelect={value => {
                    actions.editCM.selectNamespace(value);
                    actions.namespace.selectNamespace(value);
                  }}
                />
              </FormItem>
              <FormItem label="">
                <div className="form-unit">
                  <div
                    className="tc-15-table-panel config-cont-table"
                    style={{ width: '550px', maxHeight: '400px', borderBottom: '1px solid #DDDDDD' }}
                  >
                    <div className="tc-15-table-fixed-head">
                      <table className="tc-15-table-box">
                        <colgroup>
                          <col />
                          <col style={{ width: '10%' }} />
                          <col />
                          <col style={{ width: '10%' }} />
                        </colgroup>
                        <thead>
                          <tr>
                            <th>
                              <Text parent="div" overflow>
                                {t('变量名')}
                              </Text>
                            </th>
                            <th>
                              <div />
                            </th>
                            <th>
                              <Text parent="div" overflow>
                                {t('变量值')}
                              </Text>
                            </th>
                            <th>
                              <div />
                            </th>
                          </tr>
                        </thead>
                      </table>
                    </div>
                    <div className="tc-15-table-fixed-body" style={{ maxHeight: '270px' }}>
                      <table className="tc-15-table-box tc-15-table-rowhover">
                        <colgroup>
                          <col />
                          <col style={{ width: '10%' }} />
                          <col />
                          <col style={{ width: '10%' }} />
                        </colgroup>
                        <tbody>{this._renderVariableList()}</tbody>
                      </table>
                    </div>
                  </div>

                  <LinkButton disabled={!canAdd} errorTip={t('请先完成待编辑项')} onClick={actions.editCM.addVariable}>
                    {t('添加变量')}
                  </LinkButton>
                </div>
              </FormItem>
              <li className="pure-text-row fixed">
                <Button
                  className="mr10"
                  type="primary"
                  disabled={modifyResourceFlow.operationState === OperationState.Performing}
                  onClick={this._handleSubmit.bind(this)}
                >
                  {failed ? t('重试') : t('创建ConfigMap')}
                </Button>
                <Button onClick={e => router.navigate(Object.assign({}, urlParams, { mode: 'list' }), route.queries)}>
                  {t('取消')}
                </Button>
                <TipInfo
                  isShow={failed}
                  className="error"
                  style={{ display: 'inline-block', marginLeft: '20px', marginBottom: '0px', maxWidth: '750px' }}
                >
                  {getWorkflowError(modifyResourceFlow)}
                </TipInfo>
              </li>
            </ul>
          </div>
        </FormLayout>
      </MainBodyLayout>
    );
  }

  /** 处理提交请求 */
  private _handleSubmit() {
    let { actions, subRoot, route, region } = this.props,
      { resourceInfo, mode, cmEdit } = subRoot;

    actions.validate.cm.validateCMEdit();

    if (validateCMActions._validateCMEdit(cmEdit)) {
      let { name, namespace } = cmEdit;

      let jsonData = {
        kind: 'ConfigMap',
        apiVersion: (resourceInfo.group ? resourceInfo.group + '/' : '') + resourceInfo.version,
        metadata: {
          name: name,
          namespace: reduceNs(namespace)
        },
        data: cmEdit.variables.reduce((prev, next) => {
          return Object.assign({}, prev, {
            [next.key]: next.value
          });
        }, {})
      };

      let resource: CreateResource = {
        id: uuid(),
        resourceInfo,
        namespace: cmEdit.namespace || 'default',
        mode,
        clusterId: route.queries['clusterId'],
        jsonData: JSON.stringify(jsonData)
      };

      actions.workflow.modifyResource.start([resource], region.selection.value);
      actions.workflow.modifyResource.perform();
    }
  }

  /** 创建ConfigMap 的按钮 */
  private _handleClickForCreateConfigMap(urlParams) {
    let { actions, route } = this.props;

    router.navigate(Object.assign({}, urlParams, { resourceName: 'configmap' }), route.queries);

    // 需要拉取namsepsace，并且要进行路由的跳转
    actions.resource.initResourceInfoAndFetchData(true, 'configmap');
  }
}
