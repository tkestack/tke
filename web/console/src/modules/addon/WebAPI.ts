import { RecordSet, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { apiServerVersion } from '../../../config';
import {
    Method, operationResult, reduceK8sQueryString, reduceK8sRestfulPath, reduceNetworkRequest,
    reduceNetworkWorkflow, requestMethodForAction
} from '../../../helpers';
import {
    CreateResource, RequestParams, Resource, ResourceFilter, ResourceInfo, UserDefinedHeader
} from '../common';

const tips = seajs.require('tips');
