import { initValidator } from './../../common/models/Validation';
import { extend, deepClone, uuid } from '@tencent/qcloud-lib';
import { generateFetcherActionCreator } from '@tencent/qcloud-redux-fetcher';
import { generateWorkflowActionCreator, OperationTrigger, isSuccessWorkflow } from '@tencent/qcloud-redux-workflow';
import { RootState, Image, ImageFilter, ImageCreation } from '../models';
import * as ActionType from '../constants/ActionType';
import { InitImage } from '../constants/Config';
import * as WebAPI from '../WebAPI';
import { router } from '../router';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { createListAction } from '@tencent/redux-list';

type GetState = () => RootState;

const FFModelImageActions = createListAction<Image, ImageFilter>({
  actionName: 'image',
  fetcher: async (query, getState: GetState) => {
    let response = await WebAPI.fetchImageList(query);
    return response;
  },
  getRecord: (getState: GetState) => {
    return getState().image;
  }
});

const restActions = {
  /** 创建 Image */
  createImage: generateWorkflowActionCreator<ImageCreation, void>({
    actionType: ActionType.CreateImage,
    workflowStateLocator: (state: RootState) => state.createImage,
    operationExecutor: WebAPI.createImage,
    after: {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { createImage, route } = getState();
        if (isSuccessWorkflow(createImage)) {
          dispatch(restActions.createImage.reset());
          dispatch(restActions.clearEdition());
          dispatch(imageActions.fetch());
          let urlParams = router.resolve(route);
          router.navigate(Object.assign({}, urlParams, { sub: 'repo', mode: 'detail' }), route.queries);
        }
      }
    }
  }),

  /** 删除 Image */
  deleteImage: generateWorkflowActionCreator<Image, void>({
    actionType: ActionType.DeleteImage,
    workflowStateLocator: (state: RootState) => state.deleteImage,
    operationExecutor: WebAPI.deleteImage,
    after: {
      [OperationTrigger.Done]: (dispatch, getState) => {
        let { deleteImage, route } = getState();
        if (isSuccessWorkflow(deleteImage)) {
          dispatch(restActions.deleteImage.reset());
          dispatch(imageActions.fetch());
        }
      }
    }
  }),

  /** --begin编辑action */
  inputImageDesc: (value: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateImageCreation,
        payload: Object.assign({}, getState().imageCreation, { displayName: value })
      });
    };
  },

  inputImageName: (value: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateImageCreation,
        payload: Object.assign({}, getState().imageCreation, { name: value })
      });
      dispatch(imageActions.validateImageName(value));
    };
  },

  selectImageVisibility: (value: string) => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateImageCreation,
        payload: Object.assign({}, getState().imageCreation, { visibility: value })
      });
    };
  },

  validateImageName(value: string) {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      let result = imageActions._validateImageName(value);
      dispatch({
        type: ActionType.UpdateImageCreation,
        payload: Object.assign({}, getState().imageCreation, { v_name: result })
      });
    };
  },

  _validateImageName(name: string) {
    let reg = /^[a-z]([-a-z0-9]*[a-z0-9])?$/,
      status = 0,
      message = '';

    if (!name) {
      status = 2;
      message = t('镜像名称不能为空');
    } else if (name.length > 63) {
      status = 2;
      message = t('镜像名称不能超过63个字符');
    } else if (!reg.test(name)) {
      status = 2;
      message = t('镜像名称格式不正确');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  },

  clearEdition: () => {
    return async (dispatch: Redux.Dispatch, getState: GetState) => {
      dispatch({
        type: ActionType.UpdateImageCreation,
        payload: InitImage
      });
    };
  },
  /** --end编辑action */

  fetchDockerRegUrl: generateFetcherActionCreator({
    actionType: ActionType.FetchDockerRegUrl,
    fetcher: async (getState: GetState, options, dispatch) => {
      let response = await WebAPI.fetchDockerRegUrl();
      return response;
    }
  })
};

export const imageActions = extend({}, FFModelImageActions, restActions);
