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
// import { connect } from 'react-redux';
// import { allActions } from '../actions';
// import * as React from 'react';
// import { bindActionCreators } from '@tencent/ff-redux';
// import { router } from '../router';
// import { RootProps } from './AlarmPolicyApp';
// import { Button, SearchBox, Table } from '@tea/component';
// import { Justify } from '@tea/component/justify';
// import { t, Trans } from '@tencent/tea-app/lib/i18n';
// const mapDispatchToProps = dispatch =>
//   Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

// @connect(
//   state => state,
//   mapDispatchToProps
// )
// export class AlarmPolicyActionPanel extends React.Component<RootProps, {}> {
//   /**新建告警设置 */
//   handleCreate() {
//     let { route, regionSelection, cluster } = this.props;
//     //actions.mode.changeMode("expand");
//     let rid = route.queries['rid'] || regionSelection.value + '',
//       clusterId = route.queries['clusterId'] || (cluster.selection ? cluster.selection.metadata.name : '');
//     router.navigate({ sub: 'create' }, { rid, clusterId });
//   }
//   render() {
//     let { alarmPolicy, actions } = this.props;
//     return (
//       <Table.ActionPanel>
//         <Justify
//           left={
//             <React.Fragment>
//               <Button type="primary" onClick={() => this.handleCreate()}>
//                 {/* <b className="icon-add" /> */}
//                 {t('新建')}
//               </Button>
//               <Button
//                 disabled={alarmPolicy.selections.length === 0}
//                 onClick={() => actions.workflow.deleteAlarmPolicy.start(alarmPolicy.selections)}
//               >
//                 {t('删除')}
//               </Button>
//             </React.Fragment>
//           }
//           right={
//             <SearchBox
//               value={alarmPolicy.query.keyword}
//               placeholder={t('请输入关键词搜索')}
//               onChange={actions.alarmPolicy.changeKeyword}
//               onSearch={actions.alarmPolicy.performSearch}
//               onClear={() => {
//                 actions.alarmPolicy.performSearch('');
//               }}
//             />
//           }
//         />
//       </Table.ActionPanel>
//     );
//   }
// }
