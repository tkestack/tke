import * as React from 'react';
import { connect } from 'react-redux';
import { Justify, Icon, Table, Button, SearchBox, Segment, Select } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux';
import { router } from '../../router';
import { allActions } from '../../actions';
import { RootProps } from './AppContainer';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface ChartActionState {
  scene?: string;
  projectID?: string;
}

@connect(state => state, mapDispatchToProps)
export class ChartActionPanel extends React.Component<RootProps, ChartActionState> {
  constructor(props, context) {
    super(props, context);
    let { route } = props;
    this.state = {
      scene: 'all',
      projectID: ''
    };
    this.changeScene(this.state.scene);
  }

  render() {
    const { actions, route, chartList, projectList } = this.props;
    let urlParam = router.resolve(route);

    let { scene, projectID } = this.state;
    let sceneOptions = [
      { value: 'all', text: t('所有模板') },
      { value: 'user', text: t('用户模板') },
      { value: 'project', text: t('业务模板') },
      { value: 'public', text: t('公共模板') }
    ];

    return (
      <React.Fragment>
        <Table.ActionPanel>
          <Justify
            left={
              <React.Fragment>
                <Segment
                  style={{ marginRight: '10px' }}
                  value={scene}
                  onChange={value => {
                    this.setState({ scene: value });
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
                      actions.chart.list.resetPaging();
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
    const { actions } = this.props;
    /** 拉取列表 */
    actions.chart.list.clearSelection();
    actions.chart.list.reset();
    // actions.chart.list.resetPaging();
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
