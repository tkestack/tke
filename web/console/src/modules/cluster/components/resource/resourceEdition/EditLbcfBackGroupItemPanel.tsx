import * as React from 'react';
import { connect } from 'react-redux';

import { Bubble, Button, Input, Justify, List, Select } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { allActions } from '../../../actions';
import { LbcfProtocolList, BackendTypeList, BackendType } from '../../../constants/Config';
import { RootProps } from '../../ClusterApp';

interface EditLbcfBackGroupItemPanelProps extends RootProps {
  backGroupId: string;
  backGroupmode: string;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditLbcfBackGroupItemPanel extends React.Component<EditLbcfBackGroupItemPanelProps, {}> {
  _renderPorts(ports) {
    let { actions, backGroupId, backGroupmode } = this.props;

    return (
      <React.Fragment>
        <List>
          {ports.map(port => {
            return (
              <List.Item key={port.id} className={port.v_portNumber.status === 2 && 'is-error'}>
                <Select
                  size={'s'}
                  options={LbcfProtocolList}
                  value={port.protocol}
                  onChange={value => actions.lbcf.updateLbcfBGPort(backGroupId, port.id + '', { protocol: value })}
                />
                <Bubble placement="right" content={port.v_portNumber.status === 2 ? port.v_portNumber.message : null}>
                  <Input
                    size={'s'}
                    style={{
                      marginLeft: 17
                    }}
                    value={port.portNumber}
                    onChange={value => actions.lbcf.updateLbcfBGPort(backGroupId, port.id + '', { portNumber: value })}
                    onBlur={e => actions.validate.lbcf.validatePort(backGroupId, port.id + '', e.target.value)}
                  />
                </Bubble>
                {/* <Button
                  disabled={ports.length === 1}
                  icon={'close'}
                  onClick={() => actions.lbcf.deleteLbcfBGPort(backGroupId, port.id + '')}
                /> */}
              </List.Item>
            );
          })}
        </List>
        {/* <Button type="link" onClick={() => actions.lbcf.addLbcfBGPort(backGroupId)}>
          {t('添加')}
        </Button> */}
      </React.Fragment>
    );
  }

  _renderSelector(labels, isNeed?: boolean) {
    let { actions, backGroupId } = this.props;

    return (
      <React.Fragment>
        <List>
          {labels.map(label => {
            return (
              <List.Item
                key={label.id}
                className={(label.v_value.status === 2 || label.v_value.status === 2) && 'is-error'}
              >
                <Bubble placement="right" content={label.v_key.status === 2 ? label.v_key.message : null}>
                  <Input
                    size={'s'}
                    value={label.key}
                    onChange={value => actions.lbcf.updateLbcfBGLabels(backGroupId, label.id + '', { key: value })}
                    onBlur={e =>
                      actions.validate.lbcf.validateLabelContent(backGroupId, label.id + '', {
                        key: e.target.value
                      })
                    }
                  />
                </Bubble>
                <FormPanel.InlineText style={{ margin: '0px 5px' }}>{'='}</FormPanel.InlineText>
                <Bubble placement="right" content={label.v_value.status === 2 ? label.v_value.message : null}>
                  <Input
                    size={'s'}
                    value={label.value}
                    onChange={value => actions.lbcf.updateLbcfBGLabels(backGroupId, label.id + '', { value: value })}
                    onBlur={e =>
                      actions.validate.lbcf.validateLabelContent(backGroupId, label.id + '', {
                        value: e.target.value
                      })
                    }
                  />
                </Bubble>
                <Button
                  disabled={isNeed && labels.length === 1}
                  icon={'close'}
                  onClick={() => actions.lbcf.deleteLbcfBGLabels(backGroupId, label.id + '')}
                />
              </List.Item>
            );
          })}
        </List>
        <Button type="link" onClick={() => actions.lbcf.addLbcfBGLabels(backGroupId)}>
          {t('添加')}
        </Button>
        {/* <FormPanel.InlineText style={{ margin: '0px 5px' }}>{'|'}</FormPanel.InlineText> */}
        {/* <Button type="link" onClick={() => actions.lbcf.showGameAppDialog(true)}>
          {t('引用Workload')}
        </Button> */}
      </React.Fragment>
    );
  }

  _renderStaticAddress(addresses) {
    let { actions, backGroupId } = this.props;

    return (
      <React.Fragment>
        <List>
          {addresses.map(address => {
            return (
              <List.Item
                key={address.id}
                className={(address.v_value.status === 2 || address.v_value.status === 2) && 'is-error'}
              >
                <Bubble placement="right" content={address.v_value.status === 2 ? address.v_value.message : null}>
                  <Input
                    size={'s'}
                    value={address.value}
                    onChange={value => actions.lbcf.updateLbcfBGAddress(backGroupId, address.id + '', { value: value })}
                    onBlur={e => actions.validate.lbcf.validateAddress(backGroupId, address.id + '')}
                  />
                </Bubble>
                <Button
                  disabled={addresses.length === 1}
                  icon={'close'}
                  onClick={() => actions.lbcf.deleteLbcfBGAddress(backGroupId, address.id + '')}
                />
              </List.Item>
            );
          })}
        </List>
        <Button type="link" onClick={() => actions.lbcf.addLbcfBGAddress(backGroupId)}>
          {t('添加')}
        </Button>
        {/* <FormPanel.InlineText style={{ margin: '0px 5px' }}>{'|'}</FormPanel.InlineText> */}
        {/* <Button type="link" onClick={() => actions.lbcf.showGameAppDialog(true)}>
          {t('引用Workload')}
        </Button> */}
      </React.Fragment>
    );
  }

  _renderPodByName(byName) {
    let { actions, backGroupId } = this.props;

    return (
      <React.Fragment>
        <List>
          {byName.map(name => {
            return (
              <List.Item
                key={name.id}
                className={(name.v_value.status === 2 || name.v_value.status === 2) && 'is-error'}
              >
                <Bubble placement="right" content={name.v_value.status === 2 ? name.v_value.message : null}>
                  <Input
                    size={'s'}
                    value={name.value}
                    onChange={value => actions.lbcf.updateLbcfBGPodName(backGroupId, name.id + '', { value: value })}
                    onBlur={e => actions.validate.lbcf.validatePodName(backGroupId, name.id + '')}
                  />
                </Bubble>
                <Button
                  disabled={name.length === 1}
                  icon={'close'}
                  onClick={() => actions.lbcf.deleteLbcfBGPodName(backGroupId, name.id + '')}
                />
              </List.Item>
            );
          })}
        </List>
        <Button type="link" onClick={() => actions.lbcf.addLbcfBGPodName(backGroupId)}>
          {t('添加')}
        </Button>
        {/* <FormPanel.InlineText style={{ margin: '0px 5px' }}>{'|'}</FormPanel.InlineText> */}
        {/* <Button type="link" onClick={() => actions.lbcf.showGameAppDialog(true)}>
          {t('引用Workload')}
        </Button> */}
      </React.Fragment>
    );
  }
  render() {
    let { actions, subRoot, backGroupId, backGroupmode } = this.props,
      { lbcfEdit } = subRoot,
      { lbcfBackGroupEditions } = lbcfEdit;

    let backGroupItem = lbcfBackGroupEditions.find(item => item.id === this.props.backGroupId);
    let {
      id,
      ports,
      labels,
      name,
      v_name,
      onEdit,
      serviceName,
      staticAddress,
      backgroupType,
      v_serviceName,
      byName
    } = backGroupItem;
    let canSave = v_name.status === 1;
    if (backgroupType === BackendType.Pods) {
      canSave =
        canSave &&
        ports.every(port => port.v_portNumber.status === 1) &&
        labels.every(label => label.v_key.status === 1 && label.v_value.status === 1) &&
        byName.every(name => name.v_value.status === 1);
    } else if (backgroupType === BackendType.Pods) {
      canSave =
        canSave &&
        ports.every(port => port.v_portNumber.status === 1) &&
        labels.every(label => label.v_key.status === 1 && label.v_value.status === 1) &&
        v_serviceName.status === 1;
    } else {
      canSave = canSave && staticAddress.every(address => address.v_value.status === 1);
    }

    return (
      <FormPanel
        isNeedCard={false}
        fixed
        key={id}
        fieldStyle={{
          minWidth: 400
        }}
      >
        {backGroupmode === 'create' ? (
          <FormPanel.Item
            label={t('名称')}
            errorTipsStyle="Icon"
            validator={v_name}
            input={{
              value: name,
              placeholder: t('请输入名称'),
              onChange: value => actions.lbcf.inputLbcfBackGroupName(backGroupId, value),
              onBlur: e => actions.validate.lbcf.validateLbcfBackGroupName(backGroupId, e.target.value)
            }}
            message={t('最长63个字符，只能包含小写字母、数字及分隔符("-")，且必须以小写字母开头，数字或小写字母结尾')}
          />
        ) : (
          <FormPanel.Item text label={t('名称')}>
            {name}
          </FormPanel.Item>
        )}
        {backGroupmode === 'create' ? (
          <FormPanel.Item
            label={t('类型')}
            segment={{
              value: backgroupType,
              options: BackendTypeList,
              onChange: value => actions.lbcf.inputLbcfBackGroupType(backGroupId, value)
            }}
          />
        ) : (
          <FormPanel.Item text label={t('类型')}>
            {backgroupType}
          </FormPanel.Item>
        )}
        <FormPanel.Item
          label={t('Service')}
          errorTipsStyle="Icon"
          validator={v_serviceName}
          input={{
            value: serviceName,
            placeholder: t('请输入Service名称'),
            onChange: value => actions.lbcf.inputLbcfBackGroupServiceName(backGroupId, value),
            onBlur: e => actions.validate.lbcf.validateLbcfBackGroupServiceName(backGroupId, e.target.value)
          }}
          message={t('最长63个字符，只能包含小写字母、数字及分隔符("-")，且必须以小写字母开头，数字或小写字母结尾')}
          isShow={backgroupType === BackendType.Service}
        />
        <FormPanel.Item
          label={t('地址')}
          text={staticAddress.length === 0}
          isShow={backgroupType === BackendType.Static}
        >
          {this._renderStaticAddress(staticAddress)}
        </FormPanel.Item>
        <FormPanel.Item label={t('端口:协议')} text={ports.length === 0} isShow={backgroupType !== BackendType.Static}>
          {this._renderPorts(ports)}
        </FormPanel.Item>
        <FormPanel.Item
          label={backgroupType !== BackendType.Pods ? t('绑定节点') : t('绑定Pod')}
          text={labels.length === 0}
          isShow={backgroupType !== BackendType.Static}
        >
          {this._renderSelector(labels, backgroupType !== BackendType.Service)}
        </FormPanel.Item>
        <FormPanel.Item label={t('Pod Name')} text={labels.length === 0} isShow={backgroupType === BackendType.Pods}>
          {this._renderPodByName(byName)}
        </FormPanel.Item>
        <FormPanel.Item isShow={backGroupmode === 'create'}>
          <Justify
            left={
              <React.Fragment>
                <Button
                  type="primary"
                  style={{ marginRight: 20 }}
                  disabled={!canSave}
                  onClick={() => {
                    actions.lbcf.changeBackgroupEditStatus(backGroupId, false);
                  }}
                >
                  保存
                </Button>
                <Button
                  disabled={lbcfBackGroupEditions.length === 1}
                  onClick={() => {
                    actions.lbcf.deleteLbcfBackGroup(backGroupId);
                  }}
                >
                  删除
                </Button>
              </React.Fragment>
            }
          />
        </FormPanel.Item>
      </FormPanel>
    );
  }
}
