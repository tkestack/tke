import * as React from 'react';
import { Button } from '@tea/component';
import { OperationState, isSuccessWorkflow } from '@tencent/ff-redux';
import { bindActionCreators, uuid } from '@tencent/qcloud-lib';
import { connect } from 'react-redux';
import * as classnames from 'classnames';
import { RootProps } from '../../ClusterApp';
import { allActions } from '../../../actions';
import { MainBodyLayout, FormLayout } from '../../../../common/layouts';
import { FormItem, InputField, TipInfo } from '../../../../common/components';
import { getWorkflowError, isEmpty } from '../../../../common/utils';
import { router } from '../../../router';
import { validateNamespaceActions } from '../../../actions/validateNamespaceActions';
import { NamespaceEditJSONYaml, CreateResource, SecretEditJSONYaml } from '../../../models';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditNamespacePanel extends React.Component<RootProps, {}> {
  componentWillUnmount() {
    let { actions } = this.props;
    actions.editNamespace.clearNamespaceEdit();
  }

  render() {
    let { actions, subRoot, route } = this.props,
      urlParams = router.resolve(route),
      { namespaceEdit, applyResourceFlow } = subRoot,
      { v_name, description, v_description } = namespaceEdit;

    let failed = applyResourceFlow.operationState === OperationState.Done && !isSuccessWorkflow(applyResourceFlow);

    return (
      <MainBodyLayout>
        <FormLayout>
          <div className="param-box server-update add">
            <ul className="form-list jiqun fixed-layout">
              <FormItem label={t('名称')}>
                <div className={classnames('form-unit', { 'is-error': v_name.status === 2 })}>
                  <InputField
                    type="text"
                    placeholder={t('请输入Namespace名称')}
                    tipMode="popup"
                    validator={v_name}
                    tip={t(
                      '最长63个字符，只能包含小写字母、数字及分隔符("-")，且必须以小写字母开头，数字或小写字母结尾'
                    )}
                    onChange={actions.editNamespace.inputNamespaceName}
                    onBlur={actions.validate.namespace.validateNamespaceName}
                  />
                </div>
              </FormItem>
              <FormItem label={t('描述')}>
                <div className={classnames('form-unit', { 'is-error': v_description.status === 2 })}>
                  <InputField
                    type="textarea"
                    placeholder={t('请输入描述信息，不超过1000个字符')}
                    tipMode="popup"
                    validator={v_description}
                    value={description}
                    onChange={actions.editNamespace.inputNamespaceDesp}
                    onBlur={actions.validate.namespace.validateNamespaceDesp}
                  />
                </div>
              </FormItem>

              <li className="pure-text-row fixed">
                <Button
                  className="mr10"
                  type="primary"
                  disabled={applyResourceFlow.operationState === OperationState.Performing}
                  onClick={this._handleSubmit.bind(this)}
                >
                  {failed ? t('重试') : t('创建Namespace')}
                </Button>
                <Button onClick={e => router.navigate(Object.assign({}, urlParams, { mode: 'list' }), route.queries)}>
                  {t('取消')}
                </Button>
                <TipInfo isShow={failed} type="error" isForm>
                  {getWorkflowError(applyResourceFlow)}
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
    let { actions, subRoot, route, region, clusterVersion } = this.props,
      { resourceInfo, mode, namespaceEdit, resourceDetailState } = subRoot;

    actions.validate.namespace.validateNamespaceEdit();

    if (validateNamespaceActions._validateNamespaceEdit(namespaceEdit)) {
      let { name, description } = namespaceEdit;

      let { clusterId, rid } = route.queries;

      /** 相关的描述 */
      let annotations = {};

      if (description) {
        annotations['description'] = description;
      }

      let jsonData: NamespaceEditJSONYaml = {
        kind: 'Namespace',
        apiVersion: (resourceInfo.group ? resourceInfo.group + '/' : '') + resourceInfo.version,
        metadata: {
          name,
          annotations: isEmpty(annotations) ? undefined : annotations
        }
      };

      /** 最终传过去的json的数据 */
      let finalJSON = JSON.stringify(jsonData);

      let resource: CreateResource = {
        id: uuid(),
        resourceInfo,
        mode,
        clusterId,
        jsonData: finalJSON
      };

      actions.workflow.applyResource.start([resource], region.selection.value);
      actions.workflow.applyResource.perform();
    }
  }
}
