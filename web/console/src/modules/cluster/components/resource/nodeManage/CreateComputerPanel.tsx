import { isSuccessWorkflow, OperationState } from '@tencent/ff-redux';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { t } from '@tencent/tea-app/lib/i18n';
import { Bubble, Button, ContentView } from '@tencent/tea-component';
import * as React from 'react';
import { connect } from 'react-redux';
import { resourceConfig } from '../../../../../../config';
import { uuid } from '../../../../../../lib/_util';
import { getWorkflowError, ResourceInfo, TipInfo } from '../../../../../modules/common';
import { allActions } from '../../../actions';
import { CreateResource, ICComponter } from '../../../models';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';
import { SelectICComputerPanel } from '../../clusterManage/SelectICComputerPanel';
import { ShowICComputerPanel } from '../../clusterManage/ShowICComputerPanel';
import { FormPanel } from '@tencent/ff-component';

interface CreateComputerState {
  clusterName: string;
  isAdding: boolean;
  computerList: ICComponter[];
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class CreateComputerPanel extends React.Component<RootProps, CreateComputerState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      clusterName: this.props.route.queries['clusterId'],
      isAdding: true,
      computerList: []
    };
  }

  goBack() {
    history.back();
  }

  addComputer() {
    this.setState({ isAdding: true });
  }

  onSaveComputer(computer: ICComponter) {
    let { computerList } = this.state,
      newComputerList = computerList.map(c => (c.isEditing ? computer : c));
    this.setState({ isAdding: false, computerList: newComputerList.slice(0) });
  }
  onCancelComputer() {
    let { computerList } = this.state;
    computerList.forEach(c => {
      c.isEditing = false;
    });
    this.setState({ isAdding: false, computerList: computerList.slice(0) });
  }
  onAddComputer(computer: ICComponter) {
    let { computerList } = this.state;
    computerList.push(computer);
    this.setState({ isAdding: false, computerList: computerList.slice(0) });
  }
  onEditComputer(index: number) {
    let { computerList } = this.state;
    computerList[index].isEditing = true;
    this.setState({ computerList: computerList.slice(0) });
  }
  onDeleteComputer(index: number) {
    let { computerList } = this.state;
    computerList.splice(index, 1);
    this.setState({ computerList: computerList.slice(0) });
    if (computerList.length === 0) {
      this.addComputer();
    }
  }
  render() {
    let {
      actions,
      subRoot: { computerState },
      route
    } = this.props;
    const workflow = computerState.createComputerWorkflow;
    const action = actions.workflow.createComputer;
    let machinesInfo: ResourceInfo = resourceConfig().machines;
    let { computerList, clusterName } = this.state;
    const cancel = () => {
      if (workflow.operationState === OperationState.Done) {
        action.reset();
      }

      if (workflow.operationState === OperationState.Started) {
        action.cancel();
      }
      action.reset();
      router.navigate({ sub: 'sub', mode: 'list', type: 'nodeManange', resourceName: 'node' }, route.queries);
    };

    const perform = () => {
      let createComputerData: CreateResource[] = [];
      computerList.forEach(computer => {
        let { ipList, password, username, privateKey, passPhrase, ssh } = computer;
        ipList.split(';').forEach(ip => {
          let labels = {};
          computer.labels.forEach(kv => {
            labels[kv.key] = kv.value;
          });
          if (computer.isGpu) {
            labels['nvidia-device-enable'] = 'enable';
          }
          let data = {
            kind: machinesInfo.headTitle,
            apiVersion: `${machinesInfo.group}/${machinesInfo.version}`,
            spec: {
              clusterName: clusterName,
              ip: ip,
              port: +ssh,
              password: password ? window.btoa(password) : undefined,
              username: username ? username : undefined,
              privateKey: privateKey ? window.btoa(privateKey) : undefined,
              passPhrase: passPhrase ? window.btoa(passPhrase) : undefined,
              labels: labels,
              type: 'Baremetal'
            }
          };
          data = JSON.parse(JSON.stringify(data));
          createComputerData.push({
            id: uuid(),
            resourceInfo: machinesInfo,
            mode: 'create',
            jsonData: JSON.stringify(data)
          });
        });
      });

      action.start(createComputerData);
      action.perform();
    };
    const failed = workflow.operationState === OperationState.Done && !isSuccessWorkflow(workflow);
    let hasEditing = computerList.filter(c => c.isEditing).length > 0 || this.state.isAdding;
    let canAdd = !hasEditing;

    let canSave = !hasEditing && computerList.length !== 0;

    return (
      <ContentView>
        <ContentView.Body>
          <FormPanel>
            <FormPanel.Item label={t('目标机器')}>
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
                  className="m"
                  type="primary"
                  disabled={workflow.operationState === OperationState.Performing}
                  onClick={perform}
                >
                  {failed ? t('重试') : t('提交')}
                </Button>
                <Button type="weak" onClick={cancel}>
                  取消
                </Button>
                <TipInfo style={{ margin: 0, display: 'inline-block' }} className="error" isShow={failed}>
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
