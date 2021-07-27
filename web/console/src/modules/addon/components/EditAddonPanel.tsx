import * as React from 'react';
import { connect } from 'react-redux';

import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators, FetchState, isSuccessWorkflow, OperationState, uuid } from '@tencent/ff-redux';
import { t } from '@tencent/tea-app/lib/i18n';
import { Alert, Button, Card, Col, Icon, Radio, Row, Table, TableColumn, Text } from '@tencent/tea-component';
import { radioable, scrollable, stylize } from 'tea-component/es/table/addons';

import { resourceConfig } from '../../../../config';
import {
  CreateResource,
  getWorkflowError,
  initValidator,
  Markdown,
  Resource,
  ResourceInfo,
  Validation
} from '../../common';
import { allActions } from '../actions';
import { validatorActions } from '../actions/validatorActions';
import { AddonNameEnum, AddonNameMap, AddonNameMapToGenerateName, ResourceNameMap } from '../constants/Config';
import { Addon, AddonEditPeJsonYaml, AddonEditUniversalJsonYaml, EsInfo, PeEdit } from '../models';
import { router } from '../router';
import { RootProps } from './AddonApp';
import { EditPersistentEventPanel } from './EditPersistentEventPanel';

// import { addonRules } from '../constants/ValidateConfig';

import { Base64 } from 'js-base64';

/**
 * 创建persistentEvent的yaml
 * @param options
 */
const ReducePersistentEventJsonData = (options: { resourceInfo: ResourceInfo; clusterId: string; peEdit: PeEdit }) => {
  let { resourceInfo, clusterId, peEdit } = options,
    { esAddress, indexName, esUsername, esPassword } = peEdit;

  let esInfo: EsInfo;

  // 处理es的相关数据
  const [scheme, addressInfo = ''] = esAddress.split('://');
  const [ipAddress, port] = addressInfo.split(':');
  esInfo = {
    ip: ipAddress,
    port: +port,
    scheme: scheme,
    indexName: indexName,
    user: esUsername,
    password: Base64.encode(esPassword)
  };

  const jsonData: AddonEditPeJsonYaml = {
    kind: resourceInfo.headTitle,
    apiVersion: (resourceInfo.group ? resourceInfo.group + '/' : '') + resourceInfo.version,
    metadata: {
      generateName: 'pe'
    },
    spec: {
      clusterName: clusterId,
      persistentBackEnd: {
        es: esInfo
      }
    }
  };
  return JSON.stringify(jsonData);
};

/**
 * 创建Helm、GameApp
 * @param options
 */
const ReduceUniversalJsonData = (options: { resourceInfo: ResourceInfo; clusterId: string }) => {
  const { resourceInfo, clusterId } = options;
  const jsonData: AddonEditUniversalJsonYaml = {
    kind: resourceInfo.headTitle,
    apiVersion: (resourceInfo.group ? resourceInfo.group + '/' : '') + resourceInfo.version,
    metadata: {
      generateName: AddonNameMapToGenerateName[resourceInfo.headTitle] || resourceInfo.requestType['list']
    },
    spec: {
      clusterName: clusterId
    }
  };
  return JSON.stringify(jsonData);
};

export interface ValidatorOptions {
  keyName: keyof AddonValidator;
  store?: any;
  value: any;
}

export interface AddonValidator {
  /** 扩展组件的名称 */
  addonName?: Validation;

  /** es的地址 */
  esAddress?: Validation;

  /** indexName */
  indexName?: Validation;

  /** logset */
  logset?: Validation;

  /** topic */
  logsetTopic?: Validation;
}

interface EdtiAddonPanelState {
  /** addon单选项选择项 */
  selected?: Addon;

  /** 表单校验项 */
  validator: AddonValidator;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class EditAddonPanel extends React.Component<RootProps, EdtiAddonPanelState> {
  constructor(props) {
    super(props);
    this.state = {
      selected: null,
      validator: {
        addonName: initValidator,
        esAddress: initValidator,
        indexName: initValidator,
        logset: initValidator,
        logsetTopic: initValidator
      }
    };
  }

  componentDidMount() {
    const { actions, region } = this.props;
    // 进行地域的拉取
    region.list.fetched !== true && actions.region.applyFilter({});
  }

  componentWillUnmount() {
    const { actions } = this.props;
    // 清除创建的相关内容
    actions.editAddon.clearCreateAddon();
  }

  /**
   * 更新校验结果
   * @param options
   */
  // changeValidatorState(options: ValidatorOptions) {
  //   let { keyName, store, value } = options;
  //   let result = validateValue(value, Object.assign({}, addonRules[keyName], { store }));
  //   this.setState({ validator: Object.assign({}, this.state.validator, { [keyName]: result }) });
  // }

  render() {
    let { route, cluster, addon, applyResourceFlow, modifyResourceFlow, editAddon } = this.props,
      urlParams = router.resolve(route);

    const { addonName } = editAddon;

    const { mode } = urlParams;

    const { clusterId, rid } = route.queries;

    const isShowLoadingForAddonInfo =
      addon.list.fetched !== true || addon.list.fetchState === FetchState.Fetching ? true : false;

    // 判断按钮是否能进行操作
    const failed =
      (applyResourceFlow.operationState === OperationState.Done && !isSuccessWorkflow(applyResourceFlow)) ||
      (modifyResourceFlow.operationState === OperationState.Done && !isSuccessWorkflow(modifyResourceFlow));

    return (
      <FormPanel>
        <FormPanel.Item text label={t('集群')} loading={!cluster.selection} loadingElement={<Icon type="loading" />}>
          <Text>{`${cluster.selection ? cluster.selection.spec.displayName : '-'}(${clusterId})`}</Text>
        </FormPanel.Item>

        <FormPanel.Item
          label={t('扩展组件')}
          text={isShowLoadingForAddonInfo ? true : false}
          loading={isShowLoadingForAddonInfo}
          loadingElement={<Icon type="loading" />}
        >
          {this._renderAllAddonList(addon.list.data.records, cluster.selection)}
        </FormPanel.Item>

        {addonName === AddonNameEnum.PersistentEvent && <EditPersistentEventPanel />}

        <FormPanel.Footer>
          <React.Fragment>
            <Button
              type="primary"
              disabled={
                applyResourceFlow.operationState === OperationState.Performing ||
                modifyResourceFlow.operationState === OperationState.Performing
              }
              onClick={this._handleSubmit.bind(this)}
            >
              {failed ? t('重试') : t('完成')}
            </Button>
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
                {getWorkflowError(modifyResourceFlow) || getWorkflowError(applyResourceFlow)}
              </Alert>
            ) : (
              <noscript />
            )}
          </React.Fragment>
        </FormPanel.Footer>
      </FormPanel>
    );
  }

  /** 展示当前已有的addon的列表 */
  private _renderAllAddonList(addon: Addon[], cluster: Resource) {
    let { editAddon, openAddon, actions } = this.props,
      { addonName } = editAddon;

    const isAlreadyOpened = (record: Addon) => {
      const finder = openAddon.list.data.records.find(
        item => item.spec.type.toLowerCase() === record.type.toLowerCase()
      );
      return finder ? true : false;
    };

    const columns: TableColumn<Addon>[] = [
      {
        key: 'name',
        header: t('组件'),
        render: x => {
          let content: React.ReactNode;
          const type = x.type;

          if (isAlreadyOpened(x)) {
            content = (
              <React.Fragment>
                <Text verticalAlign="middle">{type}</Text>
                <Text verticalAlign="middle">{`(${t('已开通')})`}</Text>
              </React.Fragment>
            );
          } else {
            content = <Text overflow>{type}</Text>;
          }

          return (
            <React.Fragment>
              {content}
              <Text parent="p">{AddonNameMap[type]}</Text>
            </React.Fragment>
          );
        }
      }
    ];

    let finalAddonList: Addon[] = [];

    //暂时去除不支持的e

    finalAddonList = addon.filter(
      item => item.type !== 'EniIpamd' && item.type !== 'GameApp' && item.type !== 'Prometheus'
    );

    return (
      <Row gap={0}>
        <Col span={6}>
          <Table
            recordKey="id"
            bordered
            columns={columns}
            records={finalAddonList}
            rowDisabled={record => isAlreadyOpened(record) || record.type === 'LogCollector'}
            addons={[
              stylize({
                headStyle: { display: 'none' },
                bodyStyle: { height: '100%' },
                style: { height: '100%' }
              }),

              scrollable({
                maxHeight: 600
              }),

              radioable({
                value: this.state.selected ? this.state.selected.id + '' : '',
                rowSelect: true,
                onChange: value => {
                  const finder = addon.find(item => item.id === value);
                  actions.editAddon.selectAddonName(finder.type);
                  this.setState({ selected: finder });
                },
                render: (element, { disabled }) => {
                  return disabled ? <Radio display="block" value /> : element;
                }
              })
            ]}
          />
        </Col>

        <Col span={15}>
          <Card bordered style={{ borderLeft: 0, height: '100%', overflow: 'auto' }}>
            <Card.Body>{this._renderAddonInfo(addon)}</Card.Body>
          </Card>
        </Col>
      </Row>
    );
  }

  /** 展示右边的说明信息 */
  private _renderAddonInfo(addon: Addon[]) {
    const { editAddon } = this.props,
      { addonName } = editAddon;

    let content: React.ReactNode;
    const addonInfo: Addon = addon.find(item => item.type === addonName);

    if (addonName === '') {
      content = <Text>{t('请在左侧选择一个扩展组件')}</Text>;
    } else if (addonInfo.description) {
      content = <Markdown style={{ maxHeight: 758, overflow: 'auto' }} text={addonInfo.description} />;
    } else {
      content = <Text>{t('暂无该扩展组件的相关说明')}</Text>;
    }

    return content;
  }

  /** 处理请求提交 */
  private _handleSubmit() {
    let { route, editAddon, clusterVersion, actions } = this.props,
      { addonName, peEdit } = editAddon;

    // 触发校验逻辑
    actions.validator.validateAddonEdit();

    if (validatorActions._validateAddonEdit(editAddon)) {
      const resourceName = ResourceNameMap[addonName] ? ResourceNameMap[addonName] : addonName;
      const resourceInfo: ResourceInfo = resourceConfig(clusterVersion)[resourceName];

      const { clusterId, rid } = route.queries;

      let finalJSON: string;

      if (addonName === AddonNameEnum.PersistentEvent) {
        finalJSON = ReducePersistentEventJsonData({ resourceInfo, clusterId, peEdit });
      } else {
        finalJSON = ReduceUniversalJsonData({ resourceInfo, clusterId });
      }

      const resource: CreateResource = {
        id: uuid(),
        resourceInfo,
        mode: 'create',
        clusterId,
        jsonData: finalJSON
      };

      actions.workflow.modifyResource.start([resource], +rid);
      actions.workflow.modifyResource.perform();
    }
  }
}
