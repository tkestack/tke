import * as RegionAPI from './RegionAPI';
import * as K8sResourceAPI from './K8sResourceAPI';

export const CommonAPI = Object.assign({}, RegionAPI, K8sResourceAPI);
