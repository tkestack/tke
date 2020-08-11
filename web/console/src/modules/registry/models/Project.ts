import { Identifiable } from '@tencent/ff-redux';
import { ChartInfoFilter } from './Chart';

export interface Project extends Identifiable {
  metadata: {
    name: string;
  };
  spec: {
    displayName: string;
  };
}

export interface ProjectFilter {
  chartInfoFilter?: ChartInfoFilter;
}
