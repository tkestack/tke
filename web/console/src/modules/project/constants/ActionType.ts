/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
/** 业务的相关操作 */
export const CreateProject = 'CreateProject';
export const EditProjectName = 'EditProjectName';
export const EditProjectManager = 'EditProjectManager';
export const EditProjecResourceLimit = 'EditProjecResourceLimit';
export const UpdateProjectEdition = 'UpdateProjectEdition';
export const ClearProjectEdition = 'ClearProjectEdition';
export const DeleteProject = 'DeleteProject';
export const AddExistMultiProject = 'AddExistMultiProject';
export const DeleteParentProject = 'DeleteParentProject';
export const ProjectDetail = 'ProjectDetail';
export const PlatformType = 'IsBussiness';

/**命名空间 */
export const CreateNamespace = 'CreateNamespace';
export const EditNamespaceResourceLimit = 'EditNamespaceResourceLimit';
export const UpdateNamespaceEdition = 'UpdateNamespaceEdition';
export const DeleteNamespace = 'DeleteNamespace';
export const MigrateNamesapce = 'MigrateNamesapce';

/**admin */
export const ModifyAdminstrator = 'ModifyAdminstrator';
export const FetchAdminstratorInfo = 'FetchAdminstratorInfo';

/** 用户 */
export const UserList = 'UserList';
export const FetchUserList = 'FetchUserList';
export const AddUser = 'AddUser';
export const RemoveUser = 'RemoveUser';
export const GetUser = 'GetUser';
export const UpdateUser = 'UpdateUser';
export const FetchUserByName = 'FetchUserByName';
export const UserStrategyList = 'UserStrategyList';
export const LoginUserInfo = 'LoginUserInfo';

/** 策略*/
export const StrategyList = 'StrategyList';
export const AddStrategy = 'AddStrategy';
export const RemoveStrategy = 'RemoveStrategy';
export const GetStrategy = 'GetStrategy';
export const UpdateStrategy = 'UpdateStrategy';
export const AssociateUser = 'AssociateUser';
export const GetCategories = 'GetCategories';
export const GetStrategyAssociatedUsers = 'GetStrategyAssociatedUsers';
export const RemoveAssociatedUser = 'RemoveAssociatedUser';
export const AddAssociatedUser = 'AddAssociatedUser';
export const GetPlatformCategories = 'GetPlatformCategories';

export const PolicyPlainList = 'PolicyPlainList';
export const PolicyAssociatedList = 'PolicyAssociatedList';
export const AssociatePolicy = 'AssociatePolicy';
export const DisassociatePolicy = 'DisassociatePolicy';
export const UpdatePolicyAssociation = 'UpdatePolicyAssociation';
export const UpdatePolicyEditorState = 'UpdatePolicyEditorState';
export const UpdatePolicyFilter = 'UpdatePolicyFilter';

/** role*/
export const RoleList = 'RoleList';
export const RolePlainList = 'RolePlainList';
export const RoleAssociatedList = 'RoleAssociatedList';
export const AddRole = 'AddRole';
export const UpdateRole = 'UpdateRole';
export const RemoveRole = 'RemoveRole';
export const AssociateRole = 'AssociateRole';
export const DisassociateRole = 'DisassociateRole';
export const UpdateRoleAssociation = 'UpdateRoleAssociation';
export const UpdateRoleCreationState = 'UpdateRoleCreationState';
export const UpdateRoleEditorState = 'UpdateRoleEditorState';
export const UpdateRoleFilter = 'UpdateRoleFilter';

/** group*/
export const GroupList = 'GroupList';
export const GroupPlainList = 'GroupPlainList';
export const GroupAssociatedList = 'GroupAssociatedList';
export const AddGroup = 'AddGroup';
export const UpdateGroup = 'UpdateGroup';
export const RemoveGroup = 'RemoveGroup';
export const AssociateGroup = 'AssociateGroup';
export const DisassociateGroup = 'DisassociateGroup';
export const UpdateGroupAssociation = 'UpdateGroupAssociation';
export const UpdateGroupCreationState = 'UpdateGroupCreationState';
export const UpdateGroupEditorState = 'UpdateGroupEditorState';
export const UpdateGroupFilter = 'UpdateGroupFilter';
