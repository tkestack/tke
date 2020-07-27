import * as React from 'react';
import { connect } from 'react-redux';
import { Justify, Icon } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { RootProps } from '../AppContainer';
import { FormPanel } from '@tencent/ff-component';
import { namespace } from '@config/resource/k8sConfig';
import { ProjectNamespace } from '@src/modules/application/models';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class HeaderPanel extends React.Component<RootProps, {}> {
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
    let { projectList, projectNamespaceList, actions, route } = this.props;
    let urlParam = router.resolve(route);
    const { mode } = urlParam;
    return (
      <Justify
        left={
          <React.Fragment>
            <h2>{t('应用管理')}</h2>
            <FormPanel.InlineText>{t('业务：')}</FormPanel.InlineText>
            <FormPanel.Select
              label={t('业务')}
              model={projectList}
              action={actions.project.list}
              value={projectList.selection ? projectList.selection.metadata.name : ''}
              onChange={value => {
                actions.project.list.selectProject(value);
              }}
              valueField={x => x.metadata.name}
              displayField={x => `${x.spec.displayName}`}
            ></FormPanel.Select>
            <FormPanel.InlineText>{t('命名空间：')}</FormPanel.InlineText>
            <FormPanel.Select
              label={'命名空间'}
              model={projectNamespaceList}
              action={actions.projectNamespace.list}
              value={projectNamespaceList.selection ? this.buildNamespace(projectNamespaceList.selection) : ''}
              valueField={x => this.buildNamespace(x)}
              displayField={x => `${x.spec.namespace}(${x.spec.clusterName})`}
              onChange={value => {
                const parts = value.split('/');
                actions.projectNamespace.list.selectProjectNamespace(
                  projectList.selection ? projectList.selection.metadata.name : '',
                  parts[0],
                  parts[1]
                );
              }}
            />
          </React.Fragment>
        }
      />
    );
  }
}
