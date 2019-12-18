import { initValidator } from '../../common/models/Validation';

export const APIKEY_URL = '/api/authv1/apikey/';
export const REPO_URL = '/apis/registry.tkestack.io/v1/namespaces/';

export const InitApiKey = {
  description: '',
  expire: 1,
  v_expire: initValidator,
  unit: 'h'
};

export const InitRepo = {
  displayName: '',
  name: '',
  v_name: initValidator,
  visibility: 'Public'
};

export const InitImage = {
  displayName: '',
  name: '',
  v_name: initValidator,
  visibility: 'Public'
};

export const Default_D_URL = 'registry.tkestack.com';
