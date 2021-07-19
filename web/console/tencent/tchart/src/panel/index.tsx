import { ColorTypes } from 'charts/index';
import * as React from "react";

import { setLocale } from "@tencent/tea-component/lib/i18n";

import { ChartFilterPanelProps } from "./containers/Filter";
import { ChartInstancesPanelProps } from "./containers/Instances";
import { ChartPanelProps } from "./containers/Pure";

setLocale((window as any).VERSION);
export { request } from "../tce/request";
export { TransformField } from "./helper";

export const ChartPanel = dynamicComponent("ChartPanel");
export const ChartFilterPanel = dynamicComponent("ChartFilterPanel");
export const ChartInstancesPanel = dynamicComponent("ChartInstancesPanel");
export { ColorTypes } from "charts/index";
export default ChartPanel;

/**
 * 动态引入chart component
 */
const ChartsComponents = React.lazy(() => import(/* webpackChunkName: "ChartsComponents" */ "./ChartsComponents"));
type ComponentName = typeof ChartsComponents extends React.LazyExoticComponent<infer ComponentProps>
  ? (ComponentProps extends (props: infer Props) => void
  ? (Props extends {componentName: infer ComponentName} ? ComponentName : any)
  : any)
  : any;

function dynamicComponent(componentName: ComponentName) {
  return function WrappedComponent(props: ChartInstancesPanelProps | ChartPanelProps | ChartFilterPanelProps) {
    return (
      <React.Suspense
        fallback={
          <div
            style={{
              width: "100%",
              height: "100%"
            }}
          >
            <span className="text-overflow">
              <i className="n-loading-icon"/>
            </span>
          </div>
        }
      >
        <ChartsComponents {...props} componentName={componentName}/>
      </React.Suspense>
    );
  };
}
