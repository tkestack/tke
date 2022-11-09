import React, { useEffect, useState } from 'react';
import {
  Alert,
  ExternalLink,
  Table,
  Form,
  Justify,
  Button,
  TableColumn,
  InputNumber,
  Select,
  InputAdornment
} from 'tea-component';
import { Trans, t } from '@tencent/tea-app/lib/i18n';
import {
  compareCpu,
  compareMemory,
  getObjectPropByPath,
  getUnitByMemory,
  transCpuUnit,
  transMemoryUnit
} from './utils';
import { MemoryUnitEnum, MemoryUnitSelectOptions } from './constants';
import { isNumber } from 'lodash';
import { useRequest } from 'ahooks';
import { fetchNamespaceLimitranges, modifyNamespaceLimitRange } from '@src/webApi/namespace';

interface ILimitRangePanelProps {
  clusterId: string;
  name: string;
}

export function LimitRangePanel({ clusterId, name }: ILimitRangePanelProps) {
  const [resource, setResource] = useState(null);

  const { data: limitRange, refresh } = useRequest(
    async () => {
      const rsp = await fetchNamespaceLimitranges({ clusterId, name });

      return rsp?.items?.at(0);
    },
    {
      ready: !!clusterId && !!name,
      onSuccess(data) {
        setResource(data || {});
      }
    }
  );

  function handleOk() {
    if (
      records.some(({ rule }) => rule?.cpu(resource)?.status === 'error' || rule?.memory(resource)?.status === 'error')
    ) {
      return;
    }

    modifyNamespaceLimitRange({ clusterId, name, resource, isCreate: !limitRange });
  }

  function handleResourceChange({ path, value }: { path: string[]; value: string }) {
    setResource(pre => {
      const copedResource = JSON.parse(JSON.stringify(pre));

      const lastIndex = path.length - 1;

      return path.reduce((rsPart, k, index) => {
        const nextPathKey = path[index + 1];

        if (index === lastIndex) {
          rsPart[k] = value;

          return copedResource;
        } else if (isNumber(nextPathKey)) {
          if (rsPart[k] === undefined) {
            rsPart[k] = [...new Array(nextPathKey + 1)].fill({});
          }

          return rsPart[k];
        } else {
          if (rsPart[k] === undefined) {
            rsPart[k] = {};
          }

          return rsPart[k];
        }
      }, copedResource);
    });
  }

  const columns: TableColumn[] = [
    {
      key: 'name',
      header: ''
    },

    {
      key: 'cpu',
      header: 'CPU',
      render({ cpuPath, rule }) {
        const cpu = getObjectPropByPath({ record: resource, path: cpuPath });

        const cpuNumber = transCpuUnit(cpu);

        const statusWithMessage = rule.cpu(resource);

        return (
          <Form.Item {...statusWithMessage}>
            <InputAdornment after={t('核')}>
              <InputNumber
                size="l"
                hideButton
                precision={3}
                min={0}
                value={cpuNumber}
                onChange={value => handleResourceChange({ path: cpuPath, value: `${value}` })}
              />
            </InputAdornment>
          </Form.Item>
        );
      }
    },

    {
      key: 'memory',
      header: t('内存'),
      render({ memoryPath, rule }) {
        const memory = getObjectPropByPath({ record: resource, path: memoryPath });

        let unit = memory ? getUnitByMemory(memory) : MemoryUnitEnum.Gi;

        // unit是否在可选择的范围内
        if (!MemoryUnitSelectOptions.map(({ value }) => value).includes(unit)) {
          unit = MemoryUnitEnum.Gi;
        }

        const memoryNumber = transMemoryUnit(memory, unit);

        const statusWithMessage = rule.memory(resource);

        return (
          <Form.Item {...statusWithMessage}>
            <InputAdornment
              after={
                <Select
                  size="xs"
                  value={unit}
                  options={MemoryUnitSelectOptions}
                  onChange={value => handleResourceChange({ path: memoryPath, value: `${memoryNumber}${value}` })}
                />
              }
            >
              <InputNumber
                size="l"
                hideButton
                precision={0}
                min={0}
                value={memoryNumber}
                onChange={value => handleResourceChange({ path: memoryPath, value: `${value}${unit}` })}
              />
            </InputAdornment>
          </Form.Item>
        );
      }
    }
  ];

  const records = [
    {
      name: t('默认Request'),
      cpuPath: ['spec', 'limits', 0, 'defaultRequest', 'cpu'],
      memoryPath: ['spec', 'limits', 0, 'defaultRequest', 'memory'],
      rule: {
        cpu: resource => {
          const limitValue = getObjectPropByPath({ record: resource, path: ['spec', 'limits', 0, 'default', 'cpu'] });
          const requestValue = getObjectPropByPath({
            record: resource,
            path: ['spec', 'limits', 0, 'defaultRequest', 'cpu']
          });

          if (compareCpu(requestValue, limitValue) > 0) {
            return {
              status: 'error',
              message: t('request限制不能超过limit限制')
            };
          }
        },

        memory: resource => {
          const limitValue = getObjectPropByPath({
            record: resource,
            path: ['spec', 'limits', 0, 'default', 'memory']
          });
          const requestValue = getObjectPropByPath({
            record: resource,
            path: ['spec', 'limits', 0, 'defaultRequest', 'memory']
          });

          if (compareMemory(requestValue, limitValue) > 0) {
            return {
              status: 'error',
              message: t('request限制不能超过limit限制')
            };
          }
        }
      }
    },

    {
      name: t('默认Limit'),
      cpuPath: ['spec', 'limits', 0, 'default', 'cpu'],
      memoryPath: ['spec', 'limits', 0, 'default', 'memory'],
      rule: {
        cpu: resource => {
          const limitValue = getObjectPropByPath({ record: resource, path: ['spec', 'limits', 0, 'default', 'cpu'] });
          const requestValue = getObjectPropByPath({
            record: resource,
            path: ['spec', 'limits', 0, 'defaultRequest', 'cpu']
          });

          if (compareCpu(requestValue, limitValue) > 0) {
            return {
              status: 'error',
              message: t('limit限制不能小于request限制')
            };
          }
        },

        memory: resource => {
          const limitValue = getObjectPropByPath({
            record: resource,
            path: ['spec', 'limits', 0, 'default', 'memory']
          });
          const requestValue = getObjectPropByPath({
            record: resource,
            path: ['spec', 'limits', 0, 'defaultRequest', 'memory']
          });

          if (compareMemory(requestValue, limitValue) > 0) {
            return {
              status: 'error',
              message: t('limit限制不能小于request限制')
            };
          }
        }
      }
    }
  ];

  return (
    <>
      <Alert style={{ marginTop: '10px', marginBottom: '10px' }}>
        <Trans>
          更多配置Limit Range指导，请参考 如何使用
          <ExternalLink
            href={'https://kubernetes.io/docs/tasks/administer-cluster/manage-resources/memory-constraint-namespace/'}
          >
            Limit Range
          </ExternalLink>
        </Trans>
      </Alert>

      <Justify
        left={
          <Button type="primary" onClick={handleOk}>
            {t('确定')}
          </Button>
        }
        right={<Button icon="refresh" onClick={refresh} />}
      />

      <Table columns={columns} records={records} />
    </>
  );
}
