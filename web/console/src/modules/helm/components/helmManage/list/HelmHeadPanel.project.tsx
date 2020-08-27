import * as React from 'react';
import { connect } from 'react-redux';

import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';
import { Justify, Tooltip } from '@tencent/tea-component';

import { allActions } from '../../../actions';
import { router } from '../../../router';
import { RootProps } from '../../HelmApp';

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
    const namespaceGroups = namespaceList.data.records.reduce((gr, { clusterDisplayName, clusterName }) => {
      const value = `${clusterDisplayName}(${clusterName})`;
      return { ...gr, [clusterName]: <Tooltip title={value}>{value}</Tooltip> };
    }, {});

    let namespaceOptions = namespaceList.data.records.map(item => {
      const text = `${item.clusterDisplayName}-${item.namespace}`;

      return {
        value: item.name,
        text: <Tooltip title={text}>{text}</Tooltip>,
        groupKey: item.clusterName,
        realText: text
      };
    });

    return (
      <Justify
        left={
          <React.Fragment>
            <h2>{t('Helm应用')}</h2>
            <FormPanel.InlineText>{t('业务：')}</FormPanel.InlineText>
            <FormPanel.Select
              label={t('业务')}
              options={projectListOptions}
              value={projectSelection}
              onChange={value => {
                actions.projectNamespace.selectProject(value);
              }}
            ></FormPanel.Select>
            <FormPanel.InlineText>{t('namespace：')}</FormPanel.InlineText>
            <FormPanel.Select
              type="simulate"
              searchable
              filter={(inputValue, { realText }: any) => realText.includes(inputValue)}
              appearence="button"
              label={'namespace'}
              groups={namespaceGroups}
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
