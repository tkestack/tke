/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
import { Justify, Icon, Table, Button, SearchBox, Segment, Select } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { RootProps } from '../ChartApp';
import { checkCustomVisible } from '@common/components/permission-provider';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface ChartActionState {
  scene?: string;
  projectID?: string;
}

@connect(state => state, mapDispatchToProps)
export class ActionPanel extends React.Component<RootProps, ChartActionState> {
  constructor(props, context) {
    super(props, context);
    const { route } = props;
    const urlParams = router.resolve(route);
    this.state = {
      scene: urlParams['tab'] || 'all',
      projectID: ''
    };
    this.changeScene(this.state.scene);
  }

  render() {
    const { actions, route, chartList, projectList, chartGroupList } = this.props;
    const urlParam = router.resolve(route);
    const { sub } = urlParam;

    const { scene, projectID } = this.state;
    const sceneOptions = [
      { value: 'all', text: t('所有模板') },
      { value: 'user', text: t('用户模板') },
      ...(checkCustomVisible('platform.registry.project') ? [{ value: 'project', text: t('业务模板') }] : []),
      { value: 'public', text: t('公共模板') }
    ];

    return (
      <React.Fragment>
        <Table.ActionPanel>
          <Justify
            left={
              <React.Fragment>
                {/* <Button
                  type="primary"
                  onClick={e => {
                    e.preventDefault();
                    router.navigate({ mode: 'create', sub: 'chart' }, route.queries);
                  }}
                >
                  {t('新建')}
                </Button> */}
                <Segment
                  value={scene}
                  onChange={value => {
                    this.setState({ scene: value });
                    router.navigate({ mode: 'list', sub: 'chart', tab: value }, route.queries);
                    this.changeScene(value);
                  }}
                  options={sceneOptions}
                />
                {this.state.scene === 'project' && (
                  <FormPanel.Select
                    placeholder={'请选择业务'}
                    value={projectID}
                    model={projectList}
                    action={actions.project.list}
                    valueField={x => x.metadata.name}
                    displayField={x => `${x.spec.displayName}`}
                    onChange={value => {
                      this.setState({ projectID: value });
                      /** 拉取列表 */
                      actions.chart.list.reset();
                      actions.chart.list.applyFilter({
                        repoType: this.state.scene,
                        projectID: value
                      });
                    }}
                  />
                )}
              </React.Fragment>
            }
            right={
              <React.Fragment>
                <SearchBox
                  value={chartList.query.keyword || ''}
                  onChange={actions.chart.list.changeKeyword}
                  onSearch={actions.chart.list.performSearch}
                  onClear={() => {
                    actions.chart.list.performSearch('');
                  }}
                  placeholder={t('请输入Chart名称')}
                />
              </React.Fragment>
            }
          />
        </Table.ActionPanel>
      </React.Fragment>
    );
  }

  private changeScene(scene: string) {
    const { actions, route, chartList, projectList } = this.props;
    /** 拉取列表 */
    actions.chart.list.reset();
    actions.chart.list.applyFilter({
      repoType: scene
    });
    /** 获取具备权限的业务列表 */
    if (scene === 'project') {
      actions.project.list.fetch();
      this.setState({ projectID: '' });
    }
  }
}
