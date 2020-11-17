import * as React from 'react';
import { connect } from 'react-redux';

import { Alert, Button, Bubble, ContentView, Icon, Justify, Switch, Text } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import {
    bindActionCreators, FetchState, isSuccessWorkflow, OperationState, uuid
} from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';

import { InputField } from '../../common/components';
import { getWorkflowError } from '../../common/utils';
import { allActions } from '../actions';
import { validatorActions } from '../actions/validatorActions';
import { CreateResource, EsInfo, PeEditJSONYaml } from '../models';
import { router } from '../router';
import { RootProps } from './PersistentEventApp';
import { PersistentEventDeleteDialog } from './PersistentEventDeleteDialog';

import { Base64 } from 'js-base64';
import * as WebAPI from '../../logStash/WebAPI';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditPersistentEventPanel extends React.Component<RootProps, any> {
  constructor(props) {
    super(props);

    this.esDetection = this.esDetection.bind(this);

    /** checkESStatus
     * 0: init
     * 1: need es detection
     * 2: start detecting
     * 3: result success
     * 4: result failure
     */
    this.state = {
      checkESStatus: 0
    };
  }

  async esDetection() {
    this.changeESStatus(2);

    let { peEdit } = this.props;
    let { esAddress, esUsername, esPassword } = peEdit;
    let [scheme, address] = esAddress.split('://');
    let [host, port] = address.split(':');
    let ret = await WebAPI.fetchEsDetection({
      scheme: scheme,
      host: host,
      port: port,
      user: esUsername,
      password: esPassword
    }) ? 3 : 4;

    this.changeESStatus(ret);
  }

  changeESStatus(status) {
    this.setState({
      checkESStatus: status
    });
  }

  componentWillUnmount() {
    let { actions } = this.props;
    actions.editPE.clearPeEdit();
  }

  componentDidMount() {
    let { actions, route, region, peList } = this.props,
      urlParams = router.resolve(route);

    let mode = urlParams['mode'];

    if (region.list.data.recordCount === 0) {
      // 如果没有地域列表则需要拉取
      actions.region.applyFilter({});
    }

    if (mode === 'update') {
      // 如果已经拉取完 pe列表的话，可以进行初始化
      if (peList.data.recordCount) {
        // 选择当前的persistentEvent
        let peInfo = peList.data.records.find(item => item.spec.clusterName === route.queries['clusterId']);
        actions.pe.selectPe(peInfo ? [peInfo] : []);
        // 初始化用户的信息
        actions.editPE.initPeEditInfoForUpdate(peInfo);
      }
    }
  }

  render() {
    return (
      <ContentView>
        <ContentView.Header>{this._renderEditHeader()}</ContentView.Header>
        <ContentView.Body>
          {this._renderEditContainer()}
          <PersistentEventDeleteDialog />
        </ContentView.Body>
      </ContentView>
    );
  }

  /** 编辑页面的头部 */
  private _renderEditHeader() {
    return (
      <Justify
        left={
          <React.Fragment>
            <a href="javascript:;" className="back-link" onClick={this.goBack.bind(this)}>
              <Icon type="btnback" />
              {t('返回')}
            </a>
            <h2>{t('设置事件持久化')}</h2>
          </React.Fragment>
        }
      />
    );
  }

  /** 头部返回的操作 */
  private goBack() {
    let { route } = this.props,
      urlParams = router.resolve(route);

    // 回到列表处
    let routeQueries = JSON.parse(JSON.stringify(Object.assign({}, route.queries, { clusterId: undefined })));
    let newUrlParams = JSON.parse(JSON.stringify(Object.assign({}, urlParams, { mode: undefined })));
    router.navigate(newUrlParams, routeQueries);
  }

  /** 创建的实际内容 */
  private _renderEditContainer() {
    let { peEdit, actions, route, modifyPeFlow, peList, cluster } = this.props,
      urlParams = router.resolve(route),
      { isOpen, v_esAddress, esAddress, indexName, v_indexName, esUsername, esPassword } = peEdit;

    /** 当前的mode */
    let mode = urlParams['mode'];

    /** 集群的相关信息 */
    let clusterId = cluster.selection ? cluster.selection.metadata.name : '',
      clusterName = cluster.selection ? cluster.selection.spec.displayName : '';

    /** 是否编辑页面中加载数据中 */
    let isUpdateLoading = (mode === 'update' && peList.fetched !== true) || peList.fetchState === FetchState.Fetching;

    let isClusterHasCreatePE = false;

    if (peList.data.recordCount) {
      isClusterHasCreatePE = peList.data.records.find(item => item.spec.clusterName === route.queries['clusterId'])
        ? true
        : false;
    }

    let failed = modifyPeFlow.operationState === OperationState.Done && !isSuccessWorkflow(modifyPeFlow);

    const checkESStatus = this.state.checkESStatus;
    const esStatusMsg = {
      2: {
        color: 'primary',
        text: '连接中...'
      },
      3: {
        color: 'success',
        text: '连接成功！点击下方【完成】以设置事件持久化'
      },
      4: {
        color: 'warning',
        text: '连接失败！请检查 ElasticSearch 相关配置，注意开启了用户验证的 ElasticSearch 需要输入用户名和密码'
      }
    };
    const esStatusMsgColor = esStatusMsg[checkESStatus] ? esStatusMsg[checkESStatus].color : 'text';
    const esStatusMsgText = esStatusMsg[checkESStatus] ? esStatusMsg[checkESStatus].text : '';
    const esNotOK = checkESStatus === 1 || checkESStatus === 2 || checkESStatus === 4;

    return (
      <FormPanel>
        <FormPanel.Item label={t('集群')} text={true} loading={cluster.list.fetchState === FetchState.Fetching}>
          {`${clusterName}(${clusterId})`}
        </FormPanel.Item>

        <FormPanel.Item
          label={t('事件持久化存储')}
          text
          message={
            <React.Fragment>
              <Text verticalAlign="middle">{t('开启事件持久化存储功能会额外占用您集群资源 ')}</Text>
              <Text verticalAlign="middle" theme="warning">
                {t('CPU（0.2核）内存（100MB）')}
              </Text>
              <Text verticalAlign="middle">{t('。关闭本功能会释放占用的资源。')}</Text>
            </React.Fragment>
          }
        >
          <Switch
            value={isOpen}
            onChange={value => {
              actions.editPE.isOpenPE(value);
            }}
          />
        </FormPanel.Item>

        <FormPanel.Item required label={t('Elasticsearch地址')} loading={isUpdateLoading}>
          <InputField
            type="text"
            placeholder="eg: http://190.0.0.1:9200"
            tipMode="popup"
            validator={v_esAddress}
            value={esAddress}
            onChange={actions.editPE.inputEsAddress}
            onBlur={actions.validate.validateEsAddress}
          />
        </FormPanel.Item>

        <FormPanel.Item required label={t('索引')} loading={isUpdateLoading}>
          <InputField
            type="text"
            placeholder="eg: fluentd"
            tipMode="popup"
            tip={t('最长60个字符，只能包含小写字母、数字及分隔符("-"、"_"、"+")，且必须以小写字母开头')}
            validator={v_indexName}
            value={indexName}
            onChange={actions.editPE.inputIndexName}
            onBlur={actions.validate.validateIndexName}
          />
        </FormPanel.Item>

        <FormPanel.Item label={t('用户名')} loading={isUpdateLoading}>
          <InputField
            type="text"
            style={{
              width: '300px'
            }}
            placeholder="仅需要用户验证的 Elasticsearch 需要填入用户名"
            value={esUsername}
            onChange={actions.editPE.inputEsUsername}
          />
        </FormPanel.Item>

        <FormPanel.Item label={t('密码')} loading={isUpdateLoading}>
          <InputField
            type="password"
            style={{
              width: '300px'
            }}
            placeholder="仅需要用户验证的 Elasticsearch 需要填入密码"
            value={esPassword}
            onChange={actions.editPE.inputEsPassword}
          />
        </FormPanel.Item>

        <FormPanel.Item message={<Text theme={esStatusMsgColor}>{esStatusMsgText}</Text>}>
          <Justify
            left={
              <React.Fragment>
                <Button
                  disabled={!esAddress.length || !indexName.length}
                  type="primary"
                  style={{ marginRight: 20 }}
                  onClick={() => {
                    this.esDetection();
                  }}
                >
                  检测连接
                </Button>
              </React.Fragment>
            }
          />
        </FormPanel.Item>

        <FormPanel.Footer>
          <React.Fragment>
            <Bubble placement="top-start" content={!failed && esNotOK ? t('请检测连接 ElasticSearch，连接成功才能设置事件持久化') : null}>
              <Button
                type="primary"
                disabled={modifyPeFlow.operationState === OperationState.Performing || (!isOpen && !isClusterHasCreatePE) || esNotOK}
                onClick={this._handleSubmit.bind(this)}
              >
                {failed ? t('重试') : t('完成')}
              </Button>
            </Bubble>
            <Button
              type="weak"
              onClick={() => {
                router.navigate({}, route.queries);
              }}
            >
              {t('取消')}
            </Button>
            {failed ? (
              <Alert
                type="error"
                style={{ display: 'inline-block', marginLeft: '20px', marginBottom: '0px', maxWidth: '750px' }}
              >
                {getWorkflowError(modifyPeFlow)}
              </Alert>
            ) : (
              <noscript />
            )}
          </React.Fragment>
        </FormPanel.Footer>
      </FormPanel>
    );
  }

  /** 提交的校验 */
  private _handleSubmit() {
    let { actions, peEdit, route, peSelection, resourceInfo, peList, cluster } = this.props,
      { isOpen } = peEdit;

    if (!isOpen) {
      // 如果不是开启，则需要检查当前集群是否已经开启了PE，如果没有开启的话，就什么都不做
      let isClusterHasCreatePE = peList.data.records.find(
        item => item.spec.clusterName === cluster.selection.metadata.name
      )
        ? true
        : false;

      if (isClusterHasCreatePE) {
        // 如果已经创建了, 则进行删除操作，这时候需要弹窗提示
        actions.workflow.deletePeFlow.start([]);
      }
    } else {
      let { esAddress, indexName, esUsername, esPassword } = peEdit,
        urlParams = router.resolve(route);

      let regionId = route.queries['rid'];

      actions.validate.validatePeEdit();

      if (validatorActions._validatePeEdit(peEdit)) {
        let mode = urlParams['mode'];
        let isCreate = mode === 'create';

        let clusterName = route.queries['clusterId'];

        // 处理es的相关数据
        let [scheme, addressInfo = ''] = esAddress.split('://');
        let [ipAddress, port] = addressInfo.split(':');
        let esInfo: EsInfo = {
          ip: ipAddress,
          port: +port,
          scheme: scheme,
          indexName: indexName,
          user: esUsername,
          password: Base64.encode(esPassword)
        };

        let jsonData: PeEditJSONYaml;

        if (isCreate) {
          jsonData = {
            kind: resourceInfo.headTitle,
            apiVersion: (resourceInfo.group ? resourceInfo.group + '/' : '') + resourceInfo.version,
            metadata: {
              generateName: 'pe'
            },
            spec: {
              clusterName,
              persistentBackEnd: {
                es: esInfo
              }
            }
          };
        } else {
          jsonData = {
            spec: {
              persistentBackEnd: {
                es: esInfo
              }
            }
          };
        }

        // 去除当中不需要的数据
        jsonData = JSON.parse(JSON.stringify(jsonData));

        let resource: CreateResource = {
          id: uuid(),
          resourceInfo,
          mode,
          clusterId: clusterName,
          jsonData: JSON.stringify(jsonData),
          resourceIns: mode === 'update' ? peSelection[0].metadata.name : ''
        };

        actions.workflow.modifyPeFlow.start([resource], +route.queries['rid']);
        actions.workflow.modifyPeFlow.perform();
      }
    }
  }
}
