import * as React from 'react';
import { connect } from 'react-redux';
import { Justify, Icon, Table, Button, SearchBox, Segment, Select } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { RootProps } from '../ChartApp';
import { Clip, GridTable, TipDialog, WorkflowDialog } from '../../../../common/components';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface ChartActionState {
  showUsageGuideline?: boolean;
  scene?: string;
  projectID?: string;
}

@connect(state => state, mapDispatchToProps)
export class ActionPanel extends React.Component<RootProps, ChartActionState> {
  constructor(props, context) {
    super(props, context);
    let { route } = props;
    let urlParams = router.resolve(route);
    this.state = {
      showUsageGuideline: false,
      scene: urlParams['tab'] || 'all',
      projectID: ''
    };
    this.changeScene(this.state.scene);
  }

  render() {
    const { actions, route, chartList, projectList, chartGroupList } = this.props;
    let urlParam = router.resolve(route);
    const { sub } = urlParam;

    let { scene, projectID } = this.state;
    let sceneOptions = [
      { value: 'all', text: t('所有模板') },
      { value: 'personal', text: t('个人模板') },
      { value: 'project', text: t('业务模板') },
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
                <Button
                  type="primary"
                  onClick={() => {
                    this.setState({ showUsageGuideline: true });
                  }}
                >
                  {t('上传指引')}
                </Button>
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
        {this._renderUsageGuideDialog()}
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

  private _renderUsageGuideDialog() {
    return (
      <TipDialog
        isShow={this.state.showUsageGuideline}
        width={680}
        caption={t('Chart 上传指引')}
        cancelAction={() => this.setState({ showUsageGuideline: false })}
        performAction={() => this.setState({ showUsageGuideline: false })}
      >
        <div className="mirroring-box" style={{ marginTop: '0px' }}>
          <ul className="mirroring-upload-list">
            <li>
              <p>
                <strong>
                  <Trans>前置条件</Trans>
                </strong>
              </p>
            </li>
            <li>
              <p>
                <Trans>
                  本地安装 Helm 客户端, 更多可查看{' '}
                  <a href="https://helm.sh/docs/intro/quickstart/" target="_blank">
                    安装 Helm
                  </a>
                  .{' '}
                </Trans>
              </p>
              <code>
                <Clip target="#installHelm" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="installHelm">{`$ curl https://raw.githubusercontent.com/helm/helm/master/scripts/get | sh`}</p>
              </code>
            </li>
            <li>
              <p>
                <Trans>本地 Helm 客户端添加 TKEStack 的 repo.</Trans>
              </p>
              <code>
                <Clip target="#addTkeRepo" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="addTkeRepo">{`helm repo add [仓库名] http://${this.props.dockerRegistryUrl.data}/chart/[模板仓库名] --username tkestack --password [访问凭证] `}</p>
              </code>
              <p className="text-weak">
                <Trans>
                  获取有效访问凭证信息，请前往
                  <a
                    href="javascript:;"
                    onClick={() => {
                      let urlParams = router.resolve(this.props.route);
                      router.navigate(Object.assign({}, urlParams, { sub: 'apikey', mode: '', tab: '' }), {});
                    }}
                  >
                    [访问凭证]
                  </a>
                  管理。
                </Trans>
              </p>
            </li>
            <li>
              <p>
                <Trans>安装 helm-push 插件</Trans>
              </p>
              <code>
                <Clip target="#installHelmPush" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="installHelmPush">{`$ helm plugin install https://github.com/chartmuseum/helm-push`}</p>
              </code>
            </li>
            <li>
              <p>
                <strong>
                  <Trans>上传Helm Chart</Trans>
                </strong>
              </p>
            </li>
            <li>
              <p>
                <Trans>上传文件夹</Trans>
              </p>
              <code>
                <Clip target="#pushHelmDir" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="pushHelmDir">{`$ helm push ./myapp [仓库名]`}</p>
              </code>
            </li>
            <li>
              <p>
                <Trans>上传压缩包</Trans>
              </p>
              <code>
                <Clip target="#pushHelmTar" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="pushHelmTar">{`$ helm push myapp-1.0.1.tgz [仓库名]`}</p>
              </code>
            </li>
            <li>
              <p>
                <Trans>下载最新版本</Trans>
              </p>
              <code>
                <Clip target="#downloadChart" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="downloadChart">{`$ helm fetch [仓库名]/myapp`}</p>
              </code>
            </li>
            <li>
              <p>
                <Trans>下载指定版本</Trans>
              </p>
              <code>
                <Clip target="#downloadSChart" className="copy-btn">
                  <Trans>复制</Trans>
                </Clip>
                <p id="downloadSChart">{`$ helm fetch [仓库名]/myapp --version 1.0.1`}</p>
              </code>
            </li>
          </ul>
        </div>
      </TipDialog>
    );
  }
}
