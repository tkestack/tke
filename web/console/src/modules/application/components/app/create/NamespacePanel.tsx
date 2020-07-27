import * as React from 'react';
import { connect } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators, insertCSS, OperationState, isSuccessWorkflow } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { RootProps } from '../AppContainer';
import { router } from '../../../router';
import { FormPanel } from '@tencent/ff-component';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class NamespacePanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    const { actions, appCreation } = this.props;
    /** 拉取集群列表 */
    //不要保存filter旧数据
    actions.cluster.list.reset();
    actions.cluster.list.applyFilter();
  }

  render() {
    let { actions, route, appCreation, clusterList, namespaceList } = this.props;
    let action = actions.app.create.addAppWorkflow;

    return (
      <React.Fragment>
        <FormPanel.Item
          label={t('运行集群')}
          vkey="spec.targetCluster"
          select={{
            showRefreshBtn: true,
            value: appCreation.spec ? appCreation.spec.targetCluster : '',
            model: clusterList,
            action: actions.cluster.list,
            valueField: x => x.metadata.name,
            displayField: x => `${x.metadata.name}(${x.spec.displayName})`,
            onChange: value => {
              actions.cluster.list.selectCluster(value);
              actions.app.create.updateCreationState({
                metadata: Object.assign({}, appCreation.metadata, {
                  namespace: ''
                }),
                spec: Object.assign({}, appCreation.spec, {
                  targetCluster: value
                })
              });
            }
          }}
        ></FormPanel.Item>
        <FormPanel.Item
          label={t('命名空间')}
          vkey="metadata.namespace"
          select={{
            showRefreshBtn: true,
            value: appCreation.metadata && appCreation.metadata.namespace ? appCreation.metadata.namespace : '',
            model: namespaceList,
            action: actions.namespace.list,
            valueField: x => x.metadata.name,
            displayField: x => `${x.metadata.name}`,
            onChange: value => {
              actions.namespace.list.selectNamespace(value);
              actions.app.create.updateCreationState({
                metadata: Object.assign({}, appCreation.metadata, {
                  namespace: value
                })
              });
            }
          }}
        ></FormPanel.Item>
      </React.Fragment>
    );
  }
}
