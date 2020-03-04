import { FormPanel } from '@tencent/ff-component';
import { isSuccessWorkflow, OperationState } from '@tencent/ff-redux';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { t } from '@tencent/tea-app/lib/i18n';
import { Bubble, Button, ContentView, Icon, Justify } from '@tencent/tea-component';
import * as React from 'react';
import { connect } from 'react-redux';
import { getWorkflowError, InputField, TipInfo } from '../../../../modules/common';
import { allActions } from '../../actions';
import { GPUTYPE } from '../../constants/Config';
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
    let { actions, createIC } = this.props,
      { computerList } = createIC;
    computerList = computerList.map(c => (c.isEditing ? computer : c));
    actions.createIC.updateComputerList(computerList.slice(0));
    this.setState({ isAdding: false });
  }
  onCancelComputer() {
    let { actions, createIC } = this.props,
      { computerList } = createIC;
    computerList.forEach(c => {
      c.isEditing = false;
    });
    actions.createIC.updateComputerList(computerList.slice(0));
    this.setState({ isAdding: false });
  }
  onAddComputer(computer: ICComponter) {
    let { actions, createIC } = this.props,
      { computerList } = createIC;
    computerList.push(computer);
    actions.createIC.updateComputerList(computerList.slice(0));
    this.setState({ isAdding: false });
  }
  onEditComputer(index: number) {
    let { actions, createIC } = this.props,
      { computerList } = createIC;
    computerList[index].isEditing = true;
    actions.createIC.updateComputerList(computerList.slice(0));
  }
  onDeleteComputer(index: number) {
    let { actions, createIC } = this.props,
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
        vip,
        v_networkDevice,
        gpu,
        gpuType
      } = createIC;

    let hasEditing = computerList.filter(c => c.isEditing).length > 0 || this.state.isAdding;
    let canAdd = !hasEditing;

    let canSave = !hasEditing && v_name.status === 1 && computerList.length !== 0 && v_networkDevice.status !== 2;

    if (vip) {
      canSave = canSave && v_vipAddress.status === 1 && v_vipPort.status === 1;
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
            <FormPanel.Item label="VIP" text>
              <FormPanel.Checkbox value={vip} onChange={actions.createIC.useVip} />
            </FormPanel.Item>

            <FormPanel.Item label="" isShow={vip}>
              <InputField
                type="text"
                value={vipAddress}
                style={{ marginRight: '5px' }}
                placeholder={t('请输入 域名 或 ip地址')}
                tipMode="popup"
                validator={v_vipAddress}
                onChange={value => actions.createIC.inputVipAddress(value)}
                onBlur={actions.validate.createIC.validateVIPServer}
              />
              <InputField
                type="text"
                value={vipPort}
                style={{ width: '100px' }}
                placeholder={t('默认6443')}
                tipMode="popup"
                validator={v_vipPort}
                onChange={value => actions.createIC.inputVipPort(value)}
                onBlur={actions.validate.createIC.validatePort}
              />
            </FormPanel.Item>

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
              minMaskCode="14"
              maxMaskCode={'24'}
              value={cidr}
              onChange={(cidr, maxNodePodNum, maxClusterServiceNum) =>
                actions.createIC.setCidr(cidr, maxClusterServiceNum, maxNodePodNum)
              }
            />

            <FormPanel.Item>
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
