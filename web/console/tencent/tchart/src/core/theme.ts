/**
 * 主题配置数据信息
 */
export namespace COLORS {
  export enum Types {
    Default = "default",
    Multi = "multi"
  }

  export const Values = {
    [Types.Default]: ["#007EFA"],
    [Types.Multi]: [
      "#007EFA",
      "#29cc85",
      "#ffbb00",
      "#ff584c",
      "#9741d9",
      "#1fc0cc",
      "#7dd936",
      "#ff9c19",
      "#e63984",
      "#655ce6",
      "#47cc50",
      "#bf30b6"
    ]
  };
}

/**
 * 折线配置
 */
export namespace LINE {
  export const Width = {
    normal: 1,
    active: 2,
  }

  export const Color = {
    active: "#0064E1"
  }
}

export namespace CHART {
  export const LabelGap = 60;

  export const SpinnerTime = 6000;
}