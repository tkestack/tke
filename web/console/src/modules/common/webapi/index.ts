import * as RegionAPI from './RegionAPI';
import * as K8sResourceAPI from './K8sResourceAPI';
import * as LogAgentAPI from './LogAgentAPI';

export const CommonAPI = Object.assign({}, RegionAPI, K8sResourceAPI, LogAgentAPI);
