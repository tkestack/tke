import * as RegionAPI from './RegionAPI';
import * as K8sResourceAPI from './K8sResourceAPI';
import * as LogAgentAPI from './LogAgentAPI';
import * as PromethusAPI from './PromethusAPI';

export const CommonAPI = Object.assign({}, RegionAPI, K8sResourceAPI, LogAgentAPI, PromethusAPI);
