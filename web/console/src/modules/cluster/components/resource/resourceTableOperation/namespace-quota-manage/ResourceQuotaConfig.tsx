import { Resource } from '@src/modules/common';
import { TableColumn, Progress, InputNumber, Select, Form, InputAdornment } from 'tea-component';
import React from 'react';
import {
  getObjectPropByPath,
  getPercentWithStatus,
  getUnitByMemory,
  transMemoryUnit,
  transCpuUnit,
  transOthers,
  compareMemory,
  compareCpu,
  compareOther
} from './utils';
import { t } from '@tencent/tea-app/lib/i18n';
import { UnitTypeEnum, DefaultUnitByType, MemoryUnitEnum, MemoryUnitSelectOptions } from './constants';

interface IUseColumnProps {
  resource: Resource;
  isEditMode: boolean;
  handleResourceChange: (params: { path: string[]; value: string }) => void;
}

interface IStatusWithMessage {
  status?: 'error' | 'success';
  message?: string;
}

type IRule = (params: {
  resource: Resource;
  usedPath: string[];
  allPath: string[];
  unitType: UnitTypeEnum;
}) => IStatusWithMessage;

interface ITableRecord {
  key?: string;
  name: string;
  usedPath: string[];
  allPath: string[];
  unitType: UnitTypeEnum;
  rules?: IRule[];
}

export function getColumnByResource({ resource, isEditMode, handleResourceChange }: IUseColumnProps) {
  const baseColumn: TableColumn<ITableRecord>[] = [
    {
      key: 'resourceType',
      header: t('资源类型'),
      render({ name }) {
        return name;
      }
    },

    {
      key: 'resourceQuota',
      header: t('配额累计使用量/总配额'),
      render({ usedPath, allPath, unitType, rules }) {
        const used = getObjectPropByPath({ path: usedPath, record: resource });
        const all = getObjectPropByPath({ path: allPath, record: resource });

        let inputAfter: React.ReactNode = null;
        let unit = null;
        let usedNumber = 0;
        let allNumber = 0;
        let precision = 0;
        let inputMin = 0;
        let progressUsed = '--';
        let progressAll = '--';

        switch (unitType) {
          case UnitTypeEnum.MEMORY:
            {
              unit = all ? getUnitByMemory(all) : MemoryUnitEnum.Gi;

              // unit是否在可选择的范围内
              if (!MemoryUnitSelectOptions.map(({ value }) => value).includes(unit)) {
                unit = MemoryUnitEnum.Gi;
              }

              inputAfter = (
                <Select
                  size="xs"
                  options={MemoryUnitSelectOptions}
                  value={unit}
                  onChange={u =>
                    handleResourceChange({
                      path: allPath,
                      value: `${allNumber}${u}`
                    })
                  }
                />
              );

              precision = 0;

              usedNumber = transMemoryUnit(used, unit);
              allNumber = transMemoryUnit(all, unit);

              inputMin = usedNumber;

              progressUsed = used === null ? '--' : usedNumber.toFixed(3);
              progressAll = all === null ? '--' : allNumber.toFixed(3);
            }
            break;
          case UnitTypeEnum.CPU:
            {
              unit = DefaultUnitByType[UnitTypeEnum.CPU];

              inputAfter = unit;

              precision = 3;

              usedNumber = transCpuUnit(used);
              allNumber = transCpuUnit(all);

              inputMin = usedNumber;

              progressUsed = used === null ? '--' : usedNumber.toFixed(3);
              progressAll = all === null ? '--' : allNumber.toFixed(3);
            }
            break;
          case UnitTypeEnum.OTHER:
            {
              unit = DefaultUnitByType[UnitTypeEnum.OTHER];
              inputAfter = unit;

              precision = 0;

              usedNumber = transOthers(used);
              allNumber = transOthers(all);

              inputMin = usedNumber;

              progressUsed = used === null ? '--' : usedNumber.toFixed(0);
              progressAll = all === null ? '--' : allNumber.toFixed(0);
            }
            break;
        }

        if (isEditMode) {
          const statusWithMessage =
            rules
              ?.map(rule => rule({ resource, usedPath, allPath, unitType }))
              ?.find(({ status }) => status === 'error') ?? {};

          return (
            <Form.Item {...statusWithMessage}>
              <InputAdornment after={inputAfter}>
                <InputNumber
                  size="l"
                  hideButton
                  precision={precision}
                  min={inputMin}
                  value={allNumber}
                  onChange={value =>
                    handleResourceChange({
                      path: allPath,
                      value: `${value}${unitType === UnitTypeEnum.MEMORY ? unit : ''}`
                    })
                  }
                />
              </InputAdornment>
            </Form.Item>
          );
        } else {
          const { percent, status } = getPercentWithStatus({ used: usedNumber, all: allNumber });

          return (
            <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'center' }}>
              <Progress percent={percent} theme={status} style={{ margin: 0 }} />
              <span style={{ marginLeft: 10 }}>{`${progressUsed} / ${progressAll} ${unit}`}</span>
            </div>
          );
        }
      }
    }
  ];

  return baseColumn;
}

// validate for request and limit
function validateForRequestAndLimit({
  requestPath,
  limitPath,
  pathType
}: {
  requestPath: string[];
  limitPath: string[];
  pathType: 'request' | 'limit';
}): IRule {
  return ({ resource, unitType }) => {
    const limitValue = getObjectPropByPath({ record: resource, path: limitPath });

    const requestValue = getObjectPropByPath({ record: resource, path: requestPath });

    const errorMessages: Record<'request' | 'limit', string> = {
      request: t('request限制不能超过limit限制'),
      limit: t('limit限制不能小于request限制')
    };

    let compare = 0;

    if (unitType === UnitTypeEnum.CPU) {
      compare = compareCpu(requestValue, limitValue);
    } else if (unitType === UnitTypeEnum.MEMORY) {
      compare = compareMemory(requestValue, limitValue);
    }

    if (compare > 0) {
      return {
        status: 'error',
        message: errorMessages[pathType]
      };
    } else {
      return {};
    }
  };
}

// validate limit、request不能小于used
const validateForUsedAndCurrent: IRule = ({ resource, usedPath, allPath, unitType }) => {
  const usedValue = getObjectPropByPath({ record: resource, path: usedPath });

  const currentValue = getObjectPropByPath({ record: resource, path: allPath });

  let compare = 0;

  switch (unitType) {
    case UnitTypeEnum.CPU:
      compare = compareCpu(usedValue, currentValue);
      break;
    case UnitTypeEnum.MEMORY:
      compare = compareMemory(usedValue, currentValue);
      break;
    case UnitTypeEnum.OTHER:
      compare = compareOther(usedValue, currentValue);
      break;
  }

  if (compare > 0) {
    return {
      status: 'error',
      message: t('当前值不应该小于已使用值')
    };
  } else {
    return {};
  }
};

export const computeResourceLimitRecords: ITableRecord[] = [
  {
    name: 'CPU Request',
    usedPath: ['status', 'used', 'requests.cpu'],
    allPath: ['spec', 'hard', 'requests.cpu'],
    unitType: UnitTypeEnum.CPU,
    rules: [
      validateForRequestAndLimit({
        requestPath: ['spec', 'hard', 'requests.cpu'],
        limitPath: ['spec', 'hard', 'limits.cpu'],
        pathType: 'request'
      }),

      validateForUsedAndCurrent
    ]
  },

  {
    name: 'CPU Limit',
    usedPath: ['status', 'used', 'limits.cpu'],
    allPath: ['spec', 'hard', 'limits.cpu'],
    unitType: UnitTypeEnum.CPU,
    rules: [
      validateForRequestAndLimit({
        requestPath: ['spec', 'hard', 'requests.cpu'],
        limitPath: ['spec', 'hard', 'limits.cpu'],
        pathType: 'limit'
      }),

      validateForUsedAndCurrent
    ]
  },

  {
    name: 'Memory Request',
    usedPath: ['status', 'used', 'requests.memory'],
    allPath: ['spec', 'hard', 'requests.memory'],
    unitType: UnitTypeEnum.MEMORY,
    rules: [
      validateForRequestAndLimit({
        requestPath: ['spec', 'hard', 'requests.memory'],
        limitPath: ['spec', 'hard', 'limits.memory'],
        pathType: 'request'
      }),

      validateForUsedAndCurrent
    ]
  },

  {
    name: 'Memory Limit',
    usedPath: ['status', 'used', 'limits.memory'],
    allPath: ['spec', 'hard', 'limits.memory'],
    unitType: UnitTypeEnum.MEMORY,
    rules: [
      validateForRequestAndLimit({
        requestPath: ['spec', 'hard', 'requests.memory'],
        limitPath: ['spec', 'hard', 'limits.memory'],
        pathType: 'limit'
      }),

      validateForUsedAndCurrent
    ]
  }
];

export const storageResourceLimitRecord: ITableRecord[] = [
  {
    name: t('存储总量'),
    usedPath: ['status', 'used', 'requests.storage'],
    allPath: ['spec', 'hard', 'requests.storage'],
    unitType: UnitTypeEnum.MEMORY,
    rules: [validateForUsedAndCurrent]
  },

  {
    name: t('PVC总量'),
    usedPath: ['status', 'used', 'persistentvolumeclaims'],
    allPath: ['spec', 'hard', 'persistentvolumeclaims'],
    unitType: UnitTypeEnum.OTHER,
    rules: [validateForUsedAndCurrent]
  }
];

export const othersResourceLimitRecord: ITableRecord[] = [
  {
    name: t('Pod总量'),
    usedPath: ['status', 'used', 'pods'],
    allPath: ['spec', 'hard', 'pods'],
    unitType: UnitTypeEnum.OTHER,
    rules: [validateForUsedAndCurrent]
  },

  {
    name: t('Service总量'),
    usedPath: ['status', 'used', 'services'],
    allPath: ['spec', 'hard', 'services'],
    unitType: UnitTypeEnum.OTHER,
    rules: [validateForUsedAndCurrent]
  },

  {
    name: t('LoadBalancer类型Service总量'),
    usedPath: ['status', 'used', 'services.loadbalancers'],
    allPath: ['spec', 'hard', 'services.loadbalancers'],
    unitType: UnitTypeEnum.OTHER,
    rules: [validateForUsedAndCurrent]
  },

  {
    name: t('NodePort类型Service总量'),
    usedPath: ['status', 'used', 'services.nodeports'],
    allPath: ['spec', 'hard', 'services.nodeports'],
    unitType: UnitTypeEnum.OTHER,
    rules: [validateForUsedAndCurrent]
  },

  {
    name: t('StatefulSet总量'),
    usedPath: ['status', 'used', 'count/statefulsets.apps'],
    allPath: ['spec', 'hard', 'count/statefulsets.apps'],
    unitType: UnitTypeEnum.OTHER,
    rules: [validateForUsedAndCurrent]
  },

  {
    name: t('Deployment总量'),
    usedPath: ['status', 'used', 'count/deployments.apps'],
    allPath: ['spec', 'hard', 'count/deployments.apps'],
    unitType: UnitTypeEnum.OTHER,
    rules: [validateForUsedAndCurrent]
  },

  {
    name: t('Job总量'),
    usedPath: ['status', 'used', 'count/jobs.batch'],
    allPath: ['spec', 'hard', 'count/jobs.batch'],
    unitType: UnitTypeEnum.OTHER,
    rules: [validateForUsedAndCurrent]
  },

  {
    name: t('CronJob总量'),
    usedPath: ['status', 'used', 'count/cronjobs.batch'],
    allPath: ['spec', 'hard', 'count/cronjobs.batch'],
    unitType: UnitTypeEnum.OTHER,
    rules: [validateForUsedAndCurrent]
  },

  {
    name: t('Secret总量'),
    usedPath: ['status', 'used', 'secrets'],
    allPath: ['spec', 'hard', 'secrets'],
    unitType: UnitTypeEnum.OTHER,
    rules: [validateForUsedAndCurrent]
  },

  {
    name: t('ConfigMap总量'),
    usedPath: ['status', 'used', 'configmaps'],
    allPath: ['spec', 'hard', 'configmaps'],
    unitType: UnitTypeEnum.OTHER,
    rules: [validateForUsedAndCurrent]
  }
];
