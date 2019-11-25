import { QueryState } from '@tencent/qcloud-redux-query';
import { ResourceInfo, RequestParams, CreateResource, UserDefinedHeader, Resource, ResourceFilter } from '../common';
import { apiServerVersion } from '../../../config';
import {
  reduceK8sRestfulPath,
  Method,
  reduceNetworkRequest,
  reduceK8sQueryString,
  operationResult,
  reduceNetworkWorkflow,
  requestMethodForAction
} from '../../../helpers';
import { RecordSet, uuid } from '@tencent/qcloud-lib';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

const tips = seajs.require('tips');
