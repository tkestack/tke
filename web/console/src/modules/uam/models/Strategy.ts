import { Identifiable } from '@tencent/ff-redux';

export interface Strategy extends Identifiable {
  spec: {
    displayName: string;
    category: string;
    description: string;
    type?: string;
    statement: {
      resources: string[];
      effect: string;
      actions: string[];
    };
  };
  [props: string]: any;
}
