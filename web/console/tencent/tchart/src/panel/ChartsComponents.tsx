import * as React from "react";
import { ChartPanel } from "./containers/Pure";
import { ChartFilterPanel } from "./containers/Filter";
import { ChartInstancesPanel } from "./containers/Instances";

const components = { ChartPanel, ChartFilterPanel, ChartInstancesPanel };
export default function(props: { componentName: keyof typeof components }) {
  return React.createElement(components[props.componentName], props);
}
