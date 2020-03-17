import { extend } from '@tencent/ff-redux';
import { RootState, UserPlain, CommonUserFilter } from '../../models';
import * as ActionTypes from '../../constants/ActionTypes';
import * as WebAPI from '../../WebAPI';
type GetState = () => RootState;

const restActions = {

};

export const listActions = extend({}, restActions);
