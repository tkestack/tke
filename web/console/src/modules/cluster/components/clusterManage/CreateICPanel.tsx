import * as React from 'react';
import { connect } from 'react-redux';

import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators, isSuccessWorkflow, OperationState } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Bubble, Button, ContentView, Icon, Justify } from '@tencent/tea-component';

import { getWorkflowError, InputField, TipInfo } from '../../../../modules/common';
import { allActions } from '../../actions';
import {
  GPUTYPE,
  CreateICVipTypeOptions,
  CreateICVipType,
  CreateICCiliumOptions,
  NetworkModeOptions
} from '../../constants/Config';
import { ICComponter } from '../../models';
import { router } from '../../router';
import { RootProps } from '../ClusterApp';
import { CIDR } from './CIDR';
import { SelectICComputerPanel } from './SelectICComputerPanel';
import { ShowICComputerPanel } from './ShowICComputerPanel';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

interface State {
  isAdding: boolean;
}

@connect(state => state, mapDispatchToProps)
export class CreateICPanel extends React.Component<RootProps, State> {
  state = {
    isAdding: true
  };
  componentDidMount() {
    this.props.actions.createIC.fetchK8sVersion();
  }
  goBack() {
    history.back();
  }
  addComputer() {
    this.setState({ isAdding: true });
  }

  componentWillUnmount() {
    this.props.actions.createIC.clear();
    this.props.actions.workflow.createIC.reset();
    this.props.actions.clusterCreation.clearClusterCreationState();
  }

  onSaveComputer(computer: ICComponter) {
    const { actions, createIC } = this.props;
    let { computerList } = createIC;
    computerList = computerList.map(c => (c.isEditing ? computer : c));
    actions.createIC.updateComputerList(computerList.slice(0));
    this.setState({ isAdding: false });
  }
  onCancelComputer() {
    const { actions, createIC } = this.props,
      { computerList } = createIC;
    computerList.forEach(c => {
      c.isEditing = false;
    });
    actions.createIC.updateComputerList(computerList.slice(0));
    this.setState({ isAdding: false });
  }
  onAddComputer(computer: ICComponter) {
    const { actions, createIC } = this.props,
      { computerList } = createIC;
    computerList.push(computer);
    actions.createIC.updateComputerList(computerList.slice(0));
    this.setState({ isAdding: false });
  }
  onEditComputer(index: number) {
    const { actions, createIC } = this.props,
      { computerList } = createIC;
    computerList[index].isEditing = true;
    actions.createIC.updateComputerList(computerList.slice(0));
  }
  onDeleteComputer(index: number) {
    const { actions, createIC } = this.props,
      { computerList } = createIC;
    computerList.splice(index, 1);
    actions.createIC.updateComputerList(computerList.slice(0));
    if (computerList.length === 0) {
      this.addComputer();
    }
  }
  render() {
    let { actions, createIC, route, createICWorkflow } = this.props,
      {
        name,
        v_name,
        k8sVersion,
        cidr,
        computerList,
        k8sVersionList,
        networkDevice,
        maxClusterServiceNum,
        maxNodePodNum,
        vipAddress,
        v_vipAddress,
        vipPort,
        v_vipPort,
        vipType,
        v_networkDevice,
        gpu,
        gpuType,
        merticsServer,
        cilium,
        networkMode,
        asNumber,
        switchIp,
        v_asNumber,
        v_switchIp,
        useBGP
      } = createIC;

    const hasEditing = computerList.filter(c => c.isEditing).length > 0 || this.state.isAdding;
    const canAdd = !hasEditing;

    let canSave =
      !hasEditing &&
      v_name.status === 1 &&
      computerList.length !== 0 &&
      v_networkDevice.status !== 2 &&
      k8sVersion !== '';

    let showExistVipUnuseTip = false;
    if (vipType === CreateICVipType.existed) {
      let nodeNum = 0;
      computerList.forEach(c => {
        nodeNum += c.ipList.split(';').length;
      });
      canSave = canSave && v_vipAddress.status === 1 && v_vipPort.status !== 2;
      showExistVipUnuseTip = nodeNum > 1 ? false : true;
    } else if (vipType === CreateICVipType.tke) {
      canSave = canSave && v_vipAddress.status === 1;
    } else if (cilium === 'Cilium' && networkMode === 'underlay' && useBGP) {
      canSave = canSave && v_asNumber.status === 1 && v_switchIp.status === 1;
    }

    const workflow = createICWorkflow;
    const action = actions.workflow.createIC;
    const cancel = () => {
      if (workflow.operationState === OperationState.Done) {
        action.reset();
      }

      if (workflow.operationState === OperationState.Started) {
        action.cancel();
      }
      action.reset();
      actions.clusterCreation.clearClusterCreationState();
      router.navigate({}, { rid: route.queries['rid'] });
    };

    const perform = () => {
      action.start([createIC]);
      action.perform();
    };
    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);

    return (
      <ContentView>
        <ContentView.Header>
          <Justify
            left={
              <React.Fragment>
                <a href="javascript:;" className="back-link" onClick={this.goBack.bind(this)}>
                  <Icon type="btnback" />
                  {t('返回')}
                </a>
                <h2>新建独立集群</h2>
              </React.Fragment>
            }
          />
        </ContentView.Header>
        <ContentView.Body>
          <FormPanel>
            <FormPanel.Item
              label={t('集群名称')}
              validator={v_name}
              input={{
                value: name,
                onChange: value => actions.createIC.inputClusterName(value),
                onBlur: actions.validate.createIC.validateClusterName,
                maxLength: 60
              }}
            />
            <FormPanel.Item
              label={t('Kubernetes版本')}
              select={{
                options: k8sVersionList,
                value: k8sVersion,
                onChange: value => actions.createIC.selectK8SVersion(value)
              }}
            />
            <FormPanel.Item
              validator={v_networkDevice}
              message={t(
                "最长63个字符，只能包含小写字母、数字及分隔符(' - ')，且必须以小写字母开头，数字或小写字母结尾"
              )}
              label={t('网卡名称')}
              input={{
                value: networkDevice,
                onChange: value => actions.createIC.inputNetworkDevice(value),
                onBlur: actions.validate.createIC.validateNetworkDevice
              }}
            />
            <FormPanel.Item label="高可用类型" text>
              <FormPanel.Segment
                value={vipType}
                options={CreateICVipTypeOptions}
                onChange={type => {
                  actions.createIC.selectVipType(type);
                  actions.createIC.inputVipAddress('');
                }}
              />
            </FormPanel.Item>

            <FormPanel.Item
              label="VIP地址"
              isShow={vipType !== CreateICVipType.unuse}
              message={
                vipType === CreateICVipType.tke ? (
                  '用户需要提供一个可用的IP地址，保证该IP和各master节点可以正常联通，TKE会为集群部署keepalived并配置该IP为VIP'
                ) : (
                  <Trans>
                    <p>
                      在用户自定义VIP情况下，VIP后端需要绑定6443（kube-apiserver端口）端口，同时请确保该VIP有至少两个LB后端（master),
                    </p>
                    <p>由于LB自身路由问题，单LB后端情况下存在集群不可用风险。</p>
                  </Trans>
                )
              }
            >
              <InputField
                type="text"
                value={vipAddress}
                style={{ marginRight: '5px' }}
                placeholder={t('请输入ip地址')}
                tipMode="popup"
                validator={v_vipAddress}
                onChange={value => actions.createIC.inputVipAddress(value)}
                onBlur={actions.validate.createIC.validateVIPServer}
              />
              {vipType === CreateICVipType.existed && (
                <React.Fragment>
                  <InputField
                    disabled
                    type="text"
                    value={vipPort}
                    style={{ width: '120px', marginRight: '5px' }}
                    tipMode="popup"
                    validator={v_vipPort}
                    onChange={value => actions.createIC.inputVipPort(value)}
                    onBlur={actions.validate.createIC.validatePort}
                  />
                  <Bubble content={t('后端6443端口的映射端口')}>
                    <Icon type="info" />
                  </Bubble>
                </React.Fragment>
              )}
            </FormPanel.Item>

            <FormPanel.Item label="mertics server" text>
              <FormPanel.Checkbox value={merticsServer} onChange={actions.createIC.useMerticsServer} />
            </FormPanel.Item>

            <FormPanel.Item
              label={t('CNI')}
              select={{
                options: CreateICCiliumOptions,
                value: cilium,
                onChange: value => actions.createIC.useCilium(value)
              }}
            />

            {cilium === 'Cilium' && (
              <FormPanel.Item
                label={t('网络模式')}
                select={{
                  options: NetworkModeOptions,
                  value: networkMode,
                  onChange: value => actions.createIC.setNetWorkMode(value)
                }}
              />
            )}

            {cilium === 'Cilium' && networkMode === 'underlay' && (
              <FormPanel.Item label="BGP" text>
                <FormPanel.Checkbox value={useBGP} onChange={actions.createIC.setUseBGP} />
              </FormPanel.Item>
            )}

            {cilium === 'Cilium' && networkMode === 'underlay' && useBGP && (
              <>
                <FormPanel.Item
                  validator={v_asNumber}
                  label={t('自治系统号')}
                  input={{
                    value: asNumber,
                    onChange: value => actions.createIC.setAsNumber(value),
                    onBlur: actions.validate.createIC.validateAsNumber
                  }}
                />

                <FormPanel.Item
                  validator={v_switchIp}
                  label={t('交换机IP')}
                  input={{
                    value: switchIp,
                    onChange: value => actions.createIC.setSwitchIp(value),
                    onBlur: actions.validate.createIC.validateSwitchIp
                  }}
                />
              </>
            )}

            <FormPanel.Item label="GPU" text>
              <FormPanel.Checkbox value={gpu} onChange={actions.createIC.useGPU} />
            </FormPanel.Item>
            <FormPanel.Item label="" isShow={gpu}>
              <FormPanel.Select
                style={{ display: 'block' }}
                value={gpuType}
                onChange={actions.createIC.inputGPUType}
                options={[
                  {
                    text: 'pGPU',
                    value: GPUTYPE.PGPU
                  },
                  {
                    text: 'vGPU',
                    value: GPUTYPE.VGPU
                  }
                ]}
              />
            </FormPanel.Item>
            <CIDR
              parts={['192', '172', '10']}
              minMaskCode="8"
              maxMaskCode={'24'}
              value={cidr}
              onChange={(cidr, maxNodePodNum, maxClusterServiceNum) =>
                actions.createIC.setCidr(cidr, maxClusterServiceNum, maxNodePodNum)
              }
            />

            <FormPanel.Item
              message={
                showExistVipUnuseTip ? (
                  <FormPanel.HelpText theme="danger">
                    {t(
                      '在用户自定义VIP情况下，请确保该VIP有至少两个LB后端（master节点），单LB后端情况下存在集群不可用风险，请增加master节点或修改【高可用类型】'
                    )}
                  </FormPanel.HelpText>
                ) : (
                  ''
                )
              }
            >
              {computerList.map((item, index) => {
                return item.isEditing ? (
                  <SelectICComputerPanel
                    key={index}
                    computer={item}
                    onSave={computer => {
                      this.onSaveComputer(computer);
                    }}
                    onCancel={() => {
                      this.onCancelComputer();
                    }}
                    isNeedGpu={true}
                  />
                ) : (
                  <ShowICComputerPanel
                    key={index}
                    computer={item}
                    canEdit={canAdd}
                    onEdit={() => {
                      this.onEditComputer(index);
                    }}
                    onDelete={() => {
                      this.onDeleteComputer(index);
                    }}
                  />
                );
              })}
              {this.state.isAdding && (
                <SelectICComputerPanel
                  onSave={computer => {
                    this.onAddComputer(computer);
                  }}
                  onCancel={() => {
                    this.onCancelComputer();
                  }}
                  isNeedGpu={true}
                />
              )}
              <div
                style={{
                  lineHeight: '44px',
                  border: '1px dashed #ddd',
                  marginTop: '10px',
                  fontSize: '12px',
                  textAlign: 'center'
                }}
              >
                <Bubble content={canAdd ? null : t('请先完成待编辑项')}>
                  <Button
                    type="link"
                    disabled={!canAdd}
                    onClick={() => {
                      this.addComputer();
                    }}
                  >
                    {t('添加')}
                  </Button>
                </Bubble>
              </div>
            </FormPanel.Item>

            <FormPanel.Footer>
              <React.Fragment>
                <Button
                  type="primary"
                  disabled={!canSave || workflow.operationState === OperationState.Performing}
                  onClick={perform}
                >
                  {failed ? t('重试') : t('提交')}
                </Button>
                <Button type="weak" onClick={cancel}>
                  {t('取消')}
                </Button>
                <TipInfo type="error" isForm isShow={failed}>
                  {getWorkflowError(workflow)}
                </TipInfo>

                <TipInfo type="error" isForm isShow={!canSave && hasEditing}>
                  {t('请先完成待编辑项')}
                </TipInfo>
              </React.Fragment>
            </FormPanel.Footer>
          </FormPanel>
        </ContentView.Body>
      </ContentView>
    );
  }
}
