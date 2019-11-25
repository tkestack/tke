import * as React from 'react';
import { ICComponter } from '../../models';
import { FormPanel, LinkButton, TipInfo } from '../../../common/components';
import { Justify, Button, Text, Radio, Segment } from '@tencent/tea-component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { authTypeList, computerRoleList } from '../../constants/Config';
import { Validation, initValidator } from '../../../common';
import { validateValue, Rule, RuleTypeEnum } from '../../../common/validate';
import { InputLabelsPanel } from './InputLabelsPanel';
import { CIDR } from './CIDR';

const rules = {
  ipList: {
    label: '目标机器',
    rules: [
      RuleTypeEnum.isRequire,
      { type: RuleTypeEnum.regExp, limit: /^\d{1,3}(\.\d{1,3}){3}(;\d{1,3}(\.\d{1,3}){3})*$/ }
    ]
  },
  ssh: {
    label: 'SSH端口',
    rules: [
      RuleTypeEnum.isRequire,
      {
        type: RuleTypeEnum.regExp,
        limit: /^\d+$/
      }
    ]
  },
  username: {
    label: '用户名',
    rules: [RuleTypeEnum.isRequire]
  },
  password: {
    label: '密码',
    rules: [RuleTypeEnum.isRequire]
  },
  privateKey: {
    label: '密钥',
    rules: [RuleTypeEnum.isRequire]
  }
};

export function SelectICComputerPanel({
  computer,
  onSave,
  onCancel,
  isNeedGpu
}: {
  computer?: ICComponter;
  onSave: (computer: ICComponter) => void;
  onCancel: () => void;
  isNeedGpu?: boolean;
}) {
  //如果使用了footer，需要在下方留出足够的空间，避免重叠

  let [ipList, setIPList] = React.useState(computer ? computer.ipList : '');
  let [ssh, setSSH] = React.useState(computer ? computer.ssh : '22');

  let [role, setRole] = React.useState(computer ? computer.role : computerRoleList[0].value);
  let [labels, setLabels] = React.useState(computer ? computer.labels : []);
  let [authType, setAuthType] = React.useState(computer ? computer.authType : authTypeList[0].value);
  let [username, setUserName] = React.useState(computer ? computer.username : '');
  let [password, setPassword] = React.useState(computer ? computer.password : '');
  let [privateKey, setPrivateKey] = React.useState(computer ? computer.privateKey : '');
  let [passPhrase, setPassPhrase] = React.useState(computer ? computer.passPhrase : '');
  let [isGpu, setIsGpu] = React.useState(computer && isNeedGpu ? computer.isGpu : false);

  let [labelsIsValid, setLabelIsValid] = React.useState(true);

  let [v_ipList, setV_ipList] = React.useState<Validation>(
    computer ? validateValue(computer.ipList, rules.ipList) : initValidator
  );
  let [v_ssh, setV_ssh] = React.useState<Validation>(
    computer ? validateValue(computer.ipList, rules.ipList) : initValidator
  );
  let [v_username, setV_username] = React.useState<Validation>(
    computer && computer.authType === 'password' ? validateValue(computer.username, rules.username) : initValidator
  );
  let [v_password, setV_password] = React.useState<Validation>(
    computer && computer.authType === 'password' ? validateValue(computer.password, rules.password) : initValidator
  );
  let [v_privateKey, setV_privateKey] = React.useState<Validation>(
    computer && computer.authType === 'cert' ? validateValue(computer.privateKey, rules.privateKey) : initValidator
  );
  let canSave = v_ipList.status === 1 && labelsIsValid;
  if (authType === 'password') {
    canSave = canSave && v_username.status === 1 && v_password.status === 1;
  } else {
    canSave = canSave && v_privateKey.status === 1;
  }

  return (
    <FormPanel
      fixed
      isNeedCard={false}
      fieldStyle={{
        minWidth: 460
      }}
    >
      <FormPanel.Item
        label={t('目标机器')}
        tips={t('建议: Master&Etcd 节点配置4核及以上的机型')}
        input={{
          value: ipList,
          onChange: setIPList,
          onBlur: () => {
            setV_ipList(validateValue(ipList, rules.ipList));
          }
        }}
        validator={v_ipList}
        errorTipsStyle="Bubble"
        message={t('可以输入多个机器IP,用“;”分隔')}
      />
      <FormPanel.Item
        label={t('SSH端口')}
        input={{
          value: ssh,
          onChange: setSSH,
          onBlur: () => {
            setV_ssh(validateValue(ssh, rules.ssh));
          }
        }}
        validator={v_ssh}
        errorTipsStyle="Bubble"
      />
      {/* <FormPanel.Item label={t('主机角色')} text>
        <Radio.Group onChange={setRole} value={role}>
          {computerRoleList.map((item, index) => (
            <Radio name={item.value} key={index}>
              {item.text}
            </Radio>
          ))}
        </Radio.Group>
      </FormPanel.Item> */}
      <FormPanel.Item
        label={t('主机label')}
        align={labels.length ? 'middle' : 'top'}
        message={t('给主机设置Label,可用于指定容器调度')}
      >
        <InputLabelsPanel
          value={labels}
          onChange={(value, isValid) => {
            setLabels(value);
            setLabelIsValid(isValid);
          }}
        />
      </FormPanel.Item>
      <FormPanel.Item label={t('认证方式')}>
        <Segment
          options={authTypeList}
          value={authType}
          onChange={value => {
            setAuthType(value);
            if (value === 'password') {
              setPassPhrase('');
              setPrivateKey('');
            } else {
              setPassword('');
              setUserName('');
            }
          }}
        />
      </FormPanel.Item>
      <FormPanel.Item
        label={t('用户名')}
        input={{
          value: username,
          onChange: setUserName,
          onBlur: () => setV_username(validateValue(username, rules.username))
        }}
        validator={v_username}
      />
      {authType === 'password' && (
        <React.Fragment>
          <FormPanel.Item
            label={t('密码')}
            input={{
              type: 'password',
              value: password,
              onChange: setPassword,
              onBlur: () => setV_password(validateValue(password, rules.password))
            }}
            validator={v_password}
          />
        </React.Fragment>
      )}
      {authType === 'cert' && (
        <React.Fragment>
          <FormPanel.Item
            label={t('私钥')}
            input={{
              multiline: true,
              style: {
                width: 400,
                resize: 'both'
              },
              value: privateKey,
              onChange: setPrivateKey,
              onBlur: () => setV_privateKey(validateValue(privateKey, rules.privateKey))
            }}
            validator={v_privateKey}
          />
          <FormPanel.Item
            label={t('私钥密码')}
            input={{
              multiline: true,
              style: {
                width: 400,
                resize: 'both'
              },
              value: passPhrase,
              onChange: setPassPhrase
            }}
          />
        </React.Fragment>
      )}
      {isNeedGpu && (
        <FormPanel.Item
          label={t('GPU')}
          message={t('使用GPU机器需提前安装驱动和runtime')}
          checkbox={{
            value: isGpu,
            onChange: setIsGpu
          }}
        />
      )}
      <FormPanel.Item>
        <Justify
          left={
            <React.Fragment>
              <Button
                type="primary"
                style={{ marginRight: 20 }}
                disabled={!canSave}
                onClick={() => {
                  onSave({
                    ipList,
                    ssh,
                    role,
                    authType,
                    username,
                    password,
                    labels,
                    privateKey,
                    passPhrase,
                    isEditing: false,
                    isGpu
                  });
                }}
              >
                保存
              </Button>
              <Button onClick={onCancel}>取消</Button>
            </React.Fragment>
          }
        />
      </FormPanel.Item>
    </FormPanel>
  );
}
