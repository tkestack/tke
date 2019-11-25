import { apiKeyActions } from './apiKeyActions';
import { repoActions } from './repoActions';
import { imageActions } from './imageActions';

export const allActions = {
  apiKey: apiKeyActions,
  repo: repoActions,
  image: imageActions
};
