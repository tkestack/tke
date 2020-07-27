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

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class HeaderPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    const { actions } = this.props;
    /** 拉取应用列表 */
    // actions.app.list.poll();
    //不要保存filter旧数据
    actions.cluster.list.reset();
    actions.cluster.list.applyFilter();
  }

  render() {
    let { clusterList, namespaceList, actions, route } = this.props;
    let urlParam = router.resolve(route);
    const { mode } = urlParam;
    return (
      <Justify
        left={
          <React.Fragment>
            <h2>{t('应用管理')}</h2>
            <FormPanel.InlineText>{t('集群：')}</FormPanel.InlineText>
            <FormPanel.Select
              label={t('集群')}
              model={clusterList}
              action={actions.cluster.list}
              value={clusterList.selection ? clusterList.selection.metadata.name : ''}
              onChange={value => {
                actions.cluster.list.selectCluster(value);
              }}
              valueField={x => x.metadata.name}
              displayField={x => `${x.spec.displayName}`}
            ></FormPanel.Select>
            <FormPanel.InlineText>{t('命名空间：')}</FormPanel.InlineText>
            <FormPanel.Select
              label={'命名空间'}
              model={namespaceList}
              action={actions.namespace.list}
              value={namespaceList.selection ? namespaceList.selection.metadata.name : ''}
              valueField={x => x.metadata.name}
              displayField={x => `${x.metadata.name}`}
              onChange={value => {
                actions.namespace.list.selectNamespace(value);
              }}
            />
          </React.Fragment>
        }
      />
    );
  }
}
