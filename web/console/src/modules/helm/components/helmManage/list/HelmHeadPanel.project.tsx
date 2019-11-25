import * as React from 'react';
import { RootProps } from '../../HelmApp';
import { RegionBar, DownMenu, DownMenuItem, FormPanel } from '../../../../common/components';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { allActions } from '../../../actions';
import { connect } from 'react-redux';

import { router } from '../../../router';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Justify, Text, Icon, Select } from '@tencent/tea-component';
import { FetchState } from '@tencent/qcloud-redux-fetcher';

const routerSea = seajs.require('router');

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class HelmHeadPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    const {
        route,
        actions,
        listState: { region, cluster }
      } = this.props,
      urlParams = router.resolve(route);
    actions.region.applyFilter({});
    actions.projectNamespace.initProjectList();
  }
  onSelect(clusterId: string) {
    let {
      actions,
      listState: { cluster }
    } = this.props;

    let item = cluster.list.data.records.find(item => item.clusterId === clusterId);
    if (item) {
      actions.cluster.select(item);
    }
  }
  render() {
    let {
      projectList,
      projectSelection,
      namespaceSelection,
      actions,
      namespaceList,
      listState: { cluster }
    } = this.props;
    let projectListOptions = projectList.map((p, index) => ({
      text: p.displayName,
      value: p.name
    }));
    let namespaceOptions = namespaceList.data.records.map((p, index) => ({
      text: `${p.name}(${cluster.selection ? cluster.selection.metadata.name : '-'})`,
      value: p.name
    }));
    return (
      <Justify
        left={
          <React.Fragment>
            <h2>{t('Helm应用')}</h2>
            <FormPanel.InlineText>{t('项目：')}</FormPanel.InlineText>
            <FormPanel.Select
              label={t('项目')}
              options={projectListOptions}
              value={projectSelection}
              onChange={value => {
                actions.projectNamespace.selectProject(value);
              }}
            ></FormPanel.Select>
            <FormPanel.InlineText>{t('namespace：')}</FormPanel.InlineText>
            <FormPanel.Select
              label={'namespace'}
              options={namespaceOptions}
              value={namespaceSelection}
              onChange={value => actions.namespace.selectNamespace(value)}
            />
          </React.Fragment>
        }
      />
    );
  }
}
