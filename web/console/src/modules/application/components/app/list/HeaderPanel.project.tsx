/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

import * as React from 'react';
import { connect } from 'react-redux';
import { Justify, Icon, Select, Tooltip } from '@tea/component';
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

    const namespaceGroups = projectNamespaceList.list.data.records.reduce((gr, { spec }) => {
      const value = `${spec.clusterDisplayName}(${spec.clusterName})`;
      return { ...gr, [spec.clusterName]: <Tooltip title={value}>{value}</Tooltip> };
    }, {});

    let namespaceOptions = projectNamespaceList.list.data.records.map(item => {
      const text = `${item.spec.clusterDisplayName}-${item.spec.namespace}`;

      return {
        value: this.buildNamespace(item),
        text: <Tooltip title={text}>{text}</Tooltip>,
        groupKey: item.spec.clusterName,
        realText: text
      };
    });

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
            <Select
              size="m"
              type="simulate"
              searchable
              filter={(inputValue, { realText }: any) => (realText ? realText.includes(inputValue) : true)}
              appearence="button"
              // label={'namespace'}
              groups={namespaceGroups}
              options={namespaceOptions}
              value={projectNamespaceList.selection ? this.buildNamespace(projectNamespaceList.selection) : ''}
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
