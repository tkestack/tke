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
import { Table, TableColumn, Text, Modal, Card, Bubble, Icon, ContentView } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { Resource } from '../../../models';
import { RootProps } from '../AppContainer';
import { YamlEditorPanel } from '../../../../common/components';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface State {
  showYamlDialog?: boolean;
  yaml?: string;
}
@connect(state => state, mapDispatchToProps)
export class ResourceTablePanel extends React.Component<RootProps, State> {
  state = {
    showYamlDialog: false,
    yaml: ''
  };
  showYaml(yaml) {
    this.setState({
      showYamlDialog: true,
      yaml
    });
  }

  _renderYamlDialog() {
    const cancel = () => this.setState({ showYamlDialog: false, yaml: '' });
    return (
      <Modal visible={true} caption={t('查看YAML')} onClose={cancel} size={700} disableEscape={true}>
        <Modal.Body>
          <YamlEditorPanel readOnly={true} config={this.state.yaml} />
        </Modal.Body>
      </Modal>
    );
  }

  /** 展示ingress的后端服务 */
  _reduceIngressRule_standalone(showData: any) {
    let httpRules = showData !== '-' ? showData : [];
    let finalRules = httpRules.map(item => {
      return {
        protocol: 'http',
        host: item.host,
        path: item.http.paths[0].path || '',
        backend: item.http.paths[0].backend
      };
    });

    const getDomain = rule => {
      return `${rule.protocol}://${rule.host}${rule.path}`;
    };

    let finalRulesLength = finalRules.length;
    return finalRules.length ? (
      <Bubble
        placement="top"
        content={finalRules.map((rule, index) => (
          <p key={index}>
            <span style={{ verticalAlign: 'middle' }}>{getDomain(finalRules[0])}</span>
            <span style={{ verticalAlign: 'middle' }}>{`-->`}</span>
            <span style={{ verticalAlign: 'middle' }}>
              {finalRules[0].backend.serviceName + ':' + finalRules[0].backend.servicePort}
            </span>
          </p>
        ))}
      >
        <p className="text-overflow" style={{ fontSize: '12px' }}>
          <span style={{ verticalAlign: 'middle' }}>{getDomain(finalRules[0])}</span>
          <span style={{ verticalAlign: 'middle' }}>{`-->`}</span>
          <span style={{ verticalAlign: 'middle' }}>
            {finalRules[0].backend.serviceName + ':' + finalRules[0].backend.servicePort}
          </span>
        </p>
        {finalRules.length > 1 && (
          <p className="text">
            <a href="javascript:;">
              <Trans count={finalRulesLength}>等{{ finalRulesLength }}条转发规则</Trans>
            </a>
          </p>
        )}
      </Bubble>
    ) : (
      <p className="text-overflow text">{t('无')}</p>
    );
  }

  render() {
    let { actions, resourceList, route } = this.props;
    const commonColumns: TableColumn<Resource>[] = [
      {
        key: 'name',
        header: t('资源名'),
        width: '20%',
        render: (x: Resource) => (
          <Text parent="div" overflow>
            {x.metadata.name || '-'}
          </Text>
        )
      },
      {
        key: 'namespace',
        header: t('命名空间'),
        render: (x: Resource) => <Text parent="div">{x.metadata.namespace || '-'}</Text>
      },
      {
        key: 'operation',
        width: '100px',
        header: t('操作'),
        render: (x: Resource) => (
          <a href="javascript:void(0)" onClick={e => this.showYaml(x.yaml)}>
            {t('查看YAML')}
          </a>
        )
      }
    ];
    let serviceColumns: TableColumn<Resource>[] = commonColumns.slice(0);
    serviceColumns.splice(
      2,
      0,
      {
        key: 'type',
        header: t('访问方式'),
        render: (x: Resource) => <Text parent="div">{(x.object && x.object.spec && x.object.spec.type) || '-'}</Text>
      },
      {
        key: 'ip',
        header: t('IP地址'),
        render: (x: Resource) => (
          <Text parent="div">{(x.object && x.object.spec && x.object.spec.clusterIP) || '-'}</Text>
        )
      }
    );
    let deploymentColumns: TableColumn<Resource>[] = commonColumns.slice(0);
    deploymentColumns.splice(2, 0, {
      key: 'replicas',
      header: t('运行/期望Pod数量'),
      render: (x: Resource) => (
        <Text parent="div">
          {(x.object &&
            x.object.spec &&
            x.object.status &&
            (x.object.status.readyReplicas || 0) + '/' + (x.object.spec.replicas || 0)) ||
            '-'}
          &nbsp;
          {(x.object &&
            x.object.spec &&
            x.object.status &&
            x.object.status.readyReplicas !== x.object.spec.replicas && <Icon type="loading" />) ||
            ''}
        </Text>
      )
    });
    let statefulsetColumns: TableColumn<Resource>[] = commonColumns.slice(0);
    statefulsetColumns.splice(2, 0, {
      key: 'replicas',
      header: t('可观察/期望Pod数量'),
      render: (x: Resource) => (
        <Text parent="div">
          {(x.object &&
            x.object.spec &&
            x.object.status &&
            (x.object.status.replicas || 0) + '/' + (x.object.spec.replicas || 0)) ||
            '-'}
          &nbsp;
          {x.object && x.object.spec && x.object.status && x.object.status.replicas !== x.object.spec.replicas && (
            <Icon type="loading" />
          )}
        </Text>
      )
    });
    let daemonsetColumns: TableColumn<Resource>[] = commonColumns.slice(0);
    daemonsetColumns.splice(2, 0, {
      key: 'replicas',
      header: t('运行/期望Pod数量'),
      render: (x: Resource) => (
        <Text parent="div">
          {(x.object &&
            x.object.status &&
            (x.object.status.numberReady || 0) + '/' + (x.object.status.desiredNumberScheduled || 0)) ||
            '-'}
          &nbsp;
          {x.object && x.object.status && x.object.status.numberReady !== x.object.status.desiredNumberScheduled && (
            <Icon type="loading" />
          )}
        </Text>
      )
    });
    let ingressColumns: TableColumn<Resource>[] = commonColumns.slice(0);
    ingressColumns.splice(2, 0, {
      key: 'backendService',
      header: t('后端服务'),
      render: (x: Resource) =>
        this._reduceIngressRule_standalone((x.object && x.object.spec && x.object.spec.rules) || '-')
    });
    let jobColumns: TableColumn<Resource>[] = commonColumns.slice(0);
    jobColumns.splice(
      2,
      0,
      {
        key: 'parallelism',
        header: t('并行度'),
        render: (x: Resource) => (
          <Text parent="div">{(x.object && x.object.spec && x.object.spec.parallelism) || '-'}</Text>
        )
      },
      {
        key: 'completions',
        header: t('重复次数'),
        render: (x: Resource) => (
          <Text parent="div">{(x.object && x.object.spec && x.object.spec.completions) || '-'}</Text>
        )
      }
    );
    let cronJobColumns: TableColumn<Resource>[] = commonColumns.slice(0);
    cronJobColumns.splice(
      2,
      0,
      {
        key: 'schedule',
        header: t('执行策略'),
        render: (x: Resource) => (
          <Text parent="div">{(x.object && x.object.spec && x.object.spec.schedule) || '-'}</Text>
        )
      },
      {
        key: 'parallelism',
        header: t('并行度'),
        render: (x: Resource) => (
          <Text parent="div">
            {(x.object && x.object.spec && x.object.spec.jobTemplate && x.object.spec.jobTemplate.spec.parallelism) ||
              '-'}
          </Text>
        )
      },
      {
        key: 'completions',
        header: t('重复次数'),
        render: (x: Resource) => (
          <Text parent="div">
            {(x.object && x.object.spec && x.object.spec.jobTemplate && x.object.spec.jobTemplate.spec.completions) ||
              '-'}
          </Text>
        )
      }
    );
    let persistentVolumeColumns: TableColumn<Resource>[] = commonColumns.slice(0);
    persistentVolumeColumns.splice(
      2,
      0,
      {
        key: 'phase',
        header: t('状态'),
        render: (x: Resource) => (
          <Text
            parent="div"
            className={
              x.object && x.object.status && x.object.status.phase && x.object.status.phase === 'Bound'
                ? 'text-success'
                : 'text-danger'
            }
          >
            {(x.object && x.object.status && x.object.status.phase) || '-'}
            &nbsp;
            {x.object && x.object.status && x.object.status.phase && x.object.status.phase !== 'Bound' && (
              <Icon type="loading" />
            )}
          </Text>
        )
      },
      {
        key: 'accessModes',
        header: t('访问权限'),
        render: (x: Resource) => (
          <Text parent="div">
            {(x.object && x.object.spec && x.object.spec.accessModes && x.object.spec.accessModes[0]) || '-'}
          </Text>
        )
      },
      {
        key: 'persistentVolumeReclaimPolicy',
        header: t('回收策略'),
        render: (x: Resource) => (
          <Text parent="div">{(x.object && x.object.spec && x.object.spec.persistentVolumeReclaimPolicy) || '-'}</Text>
        )
      },
      {
        key: 'claimRef',
        header: t('PVC'),
        render: (x: Resource) => (
          <Text parent="div">
            {(x.object && x.object.spec && x.object.spec.claimRef && x.object.spec.claimRef.name) || '-'}
          </Text>
        )
      },
      {
        key: 'storageClassName',
        header: t('StorageClass'),
        render: (x: Resource) => (
          <Text parent="div">{(x.object && x.object.spec && x.object.spec.storageClassName) || '-'}</Text>
        )
      }
    );
    let persistentVolumeClaimColumns: TableColumn<Resource>[] = commonColumns.slice(0);
    persistentVolumeClaimColumns.splice(
      2,
      0,
      {
        key: 'phase',
        header: t('状态'),
        render: (x: Resource) => (
          <Text
            parent="div"
            className={
              x.object && x.object.status && x.object.status.phase && x.object.status.phase === 'Bound'
                ? 'text-success'
                : 'text-danger'
            }
          >
            {(x.object && x.object.status && x.object.status.phase) || '-'}&nbsp;
            {x.object && x.object.status && x.object.status.phase && x.object.status.phase !== 'Bound' && (
              <Icon type="loading" />
            )}
          </Text>
        )
      },
      {
        key: 'capacity',
        header: t('Storage'),
        render: (x: Resource) => (
          <Text parent="div">
            {(x.object && x.object.status && x.object.status.capacity && x.object.status.capacity.storage) || '-'}
          </Text>
        )
      },
      {
        key: 'accessModes',
        header: t('访问权限'),
        render: (x: Resource) => (
          <Text parent="div">
            {(x.object && x.object.spec && x.object.spec.accessModes && x.object.spec.accessModes[0]) || '-'}
          </Text>
        )
      },
      {
        key: 'storageClassName',
        header: t('StorageClass'),
        render: (x: Resource) => (
          <Text parent="div">{(x.object && x.object.spec && x.object.spec.storageClassName) || '-'}</Text>
        )
      }
    );
    let storageClassColumns: TableColumn<Resource>[] = commonColumns.slice(0);
    storageClassColumns.splice(
      2,
      0,
      {
        key: 'provisioner',
        header: t('来源'),
        render: (x: Resource) => <Text parent="div">{(x.object && x.object.provisioner) || '-'}</Text>
      },
      {
        key: 'reclaimPolicy',
        header: t('回收策略'),
        render: (x: Resource) => <Text parent="div">{(x.object && x.object.reclaimPolicy) || '-'}</Text>
      }
    );
    const columnMap = {
      Service: serviceColumns,
      Deployment: deploymentColumns,
      StatefulSet: statefulsetColumns,
      Daemonset: daemonsetColumns,
      Ingress: ingressColumns,
      Job: jobColumns,
      CronJob: cronJobColumns,
      PersistentVolume: persistentVolumeColumns,
      PersistentVolumeClaim: persistentVolumeClaimColumns,
      StorageClass: storageClassColumns
    };

    let dom = [];
    resourceList.resources &&
      resourceList.resources.forEach((value, key) => {
        const kind = key.lastIndexOf('/') === -1 ? key : key.substring(key.lastIndexOf('/') + 1);
        dom.push(
          <Card key={key}>
            <Card.Body title={kind}>
              <Table
                recordKey={record => {
                  return record.id.toString();
                }}
                records={value}
                columns={columnMap[kind] || commonColumns}
              />
            </Card.Body>
          </Card>
        );
      });
    return (
      <ContentView>
        <ContentView.Body>
          {dom}
          {this.state.showYamlDialog && this._renderYamlDialog()}
        </ContentView.Body>
      </ContentView>
    );
  }
}
