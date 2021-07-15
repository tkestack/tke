import * as languages from '../i18n';

const version = (window as any).VERSION || "zh";
const language = languages[version];

export namespace TIME_PICKER {
  export const Tabs = [
      {from: "%NOW-1h", to: "%NOW", label: language.OneHour},
      {from: "%NOW-24h", to: "%NOW", label: language.OneDay},
      {from: "%NOW-168h", to: "%NOW", label: language.SevenDays},
//      {from: "%NOW-720h", to: "%NOW", label: "近30天"},
    ];
};

export namespace QUERY {
  export const Aggregation = {
    'sum': '总和',
    'count': '统计个数',
    'max': '最大值',
    'min': '最小值',
    'avg': '平均值',
  };

  export const Limit = 65535;
}

export namespace CHART {
  export const DefaultSize = {
    width: 726,
    height: 388,
  }
}

export const Kilobyte = 1024;

export const OneDayMillisecond = 86400000;