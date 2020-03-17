import { extend } from '@tencent/ff-redux';
import { RootState, Policy, PolicyFilter } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
type GetState = () => RootState;


const restActions = {

};

export const listActions = extend({}, restActions);
