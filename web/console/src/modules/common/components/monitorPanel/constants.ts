import moment from 'moment';

export interface IMonitorPanelProps {
  conditions: ICondition[];

  groups: IGroup[];

  instanceType: string;

  instanceList: string[];

  defaultSelectedInstances: string[];

  visible: boolean;
  onClose: () => void;

  title: string;
}

interface ICondition {
  key: string;
  value: string;
  expr: string;
}

interface IFields {
  expr: string;
  alias: string;
  unit: string;
}

interface IGroup {
  by: string[];
  fields: IFields[];
}

export interface IChartRenderProps {
  conditions: ICondition[];
  group: IGroup;
  timeGran: string;
  instanceType: string;
  selectedInstances: string[];
  dateRange: { startTime: moment.Moment; endTime: moment.Moment };
}

export enum DateRangeTypeEnum {
  HOUR = 'HOUR',
  DAY = 'DAY',
  WEEK = 'WEEK',
  CUSTOM = 'CUSTOM'
}

export const dateRangeTypeOptions = [
  {
    value: DateRangeTypeEnum.HOUR,
    text: '近1小时'
  },

  {
    value: DateRangeTypeEnum.DAY,
    text: '近1天'
  },

  {
    value: DateRangeTypeEnum.WEEK,
    text: '近7天'
  }
];

const _timeGranularityOptions = [
  {
    value: `timestamp(${1 * 60}s)`,
    text: '1分钟'
  },

  {
    value: `timestamp(${5 * 60}s)`,
    text: '5分钟'
  },

  {
    value: `timestamp(${60 * 60}s)`,
    text: '1小时'
  },

  {
    value: `timestamp(${24 * 60 * 60}s)`,
    text: '24小时'
  }
];

export const timeGranularityOptions = {
  [DateRangeTypeEnum.HOUR]: _timeGranularityOptions.slice(0, 2),
  [DateRangeTypeEnum.DAY]: _timeGranularityOptions.slice(0, 3),
  [DateRangeTypeEnum.WEEK]: _timeGranularityOptions.slice(2),
  [DateRangeTypeEnum.CUSTOM]: _timeGranularityOptions.slice(2)
};
