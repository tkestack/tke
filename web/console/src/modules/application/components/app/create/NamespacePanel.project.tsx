import * as React from 'react';
import { connect } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../../actions';
import { RootProps } from '../AppContainer';
import { router } from '../../../router';
import { FormPanel } from '@tencent/ff-component';
import { ProjectNamespace } from '../../../models';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class NamespacePanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    const { actions } = this.props;
    //不要保存filter旧数据
    actions.project.list.reset();
    actions.project.list.applyFilter();
  }

  buildNamespace = (x: ProjectNamespace) => {
    return x.spec.clusterName + '/' + x.spec.namespace;
  };

  splitNamespace = (x: string) => {
    return x.split('/');
  };

  render() {
    let { actions, route, appCreation, projectList, projectNamespaceList } = this.props;
    let action = actions.app.create.addAppWorkflow;
    const { appAddWorkflow } = this.props;

    return (
      <React.Fragment>
        <FormPanel.Item
          label={t('业务')}
          vkey="spec.targetCluster"
          select={{
            showRefreshBtn: true,
            value: projectList.selection ? projectList.selection.metadata.name : '',
            model: projectList,
            action: actions.project.list,
            valueField: x => x.metadata.name,
            displayField: x => `${x.spec.displayName}`,
            onChange: value => {
              actions.project.list.selectProject(value);
            }
          }}
        ></FormPanel.Item>
        <FormPanel.Item
          label={t('命名空间')}
          vkey="metadata.namespace"
          select={{
            showRefreshBtn: true,
            value: projectNamespaceList.selection ? this.buildNamespace(projectNamespaceList.selection) : '',
            model: projectNamespaceList,
            action: actions.projectNamespace.list,
            valueField: x => this.buildNamespace(x),
            displayField: x => `${x.spec.namespace}(${x.spec.clusterName})`,
            onChange: value => {
              const parts = value.split('/');
              actions.projectNamespace.list.selectProjectNamespace(
                projectList.selection.metadata.name,
                parts[0],
                parts[1]
              );

              actions.app.create.updateCreationState({
                metadata: Object.assign({}, appCreation.metadata, {
                  namespace: parts[1]
                }),
                spec: Object.assign({}, appCreation.spec, {
                  targetCluster: parts[0]
                })
              });
            }
          }}
        ></FormPanel.Item>
      </React.Fragment>
    );
  }
}
