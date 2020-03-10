import { apiKeyActions } from './apiKeyActions';
import { repoActions } from './repoActions';
import { imageActions } from './imageActions';
import { chartActions } from './chartActions';
import { chartInsActions } from './chartInsActions';

export const allActions = {
  apiKey: apiKeyActions,
  repo: repoActions,
  image: imageActions,
  chart: chartActions,
  chartIns: chartInsActions
};
