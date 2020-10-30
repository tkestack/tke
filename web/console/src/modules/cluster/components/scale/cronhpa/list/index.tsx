import React, { useState, useEffect, useContext, useCallback, useMemo } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import {
  Table,
  Justify,
  Button,
  Card,
  Layout,
  Pagination,
  Text,
  TagSearchBox,
  Modal,
  Bubble,
  Icon
} from '@tencent/tea-component';
import { RecordSet } from '@tencent/ff-redux/src';
import { router } from '@src/modules/cluster/router';
import { downloadCsv } from '@helper/downloadCsv';
import { DisplayFiledProps, Resource, ResourceInfo } from '@src/modules/common/models';
import { MetricsResourceMap } from '../constant';
import { CHANGE_NAMESPACE, StateContext, DispatchContext } from '../context';
import { isEmpty, useModal } from '@src/modules/common/utils';
import { deleteCronHpa } from '@src/modules/cluster/WebAPI/scale';
import { LinkButton, Clip, TipInfo } from '@src/modules/common/components';
import NamespaceSelect from './namespaceSelect';

const { Body, Content } = Layout;
const { autotip } = Table.addons;

interface ListProps {
  namespaces: RecordSet<Resource>;
  cronHpaData: RecordSet<Resource>;
  triggerRefresh: () => void;
  isCronHpaInstalled: boolean;
}

const List = React.memo((props: ListProps) => {
  const route = useSelector(state => state.route);
  const urlParams = router.resolve(route);
  const { clusterId, projectName } = route.queries;
  const { namespaces, cronHpaData, triggerRefresh, isCronHpaInstalled } = props;
  const { namespaceValue } = useContext(StateContext);
  const { isShowing, toggle } = useModal();
  const [removeCronHpaName, setRemoveCronHpaName] = useState();

  /**
   * 列表内容处理
   */
  const [loading, setLoading] = useState(true);
  const [emptyText, setEmptyText] = useState('');
  useEffect(() => {
    if (isEmpty(cronHpaData)) {
      setLoading(true);
      setEmptyText('');
    } else if (!cronHpaData.recordCount) {
      setLoading(false);
      setEmptyText('您选择的该资源的列表为空，您可以切换到其他命名空间');
    } else if (cronHpaData.recordCount) {
      setLoading(false);
      setEmptyText('');
    }
  }, [cronHpaData]);

  /**
   * TagSearchBox和list数据的管理
   */
  const attributes = [
    {
      type: 'input',
      key: 'resourceName',
      name: t('名称')
    }
  ];
  const [cronHpaList, setCronHpaList] = useState([]);
  const [tags, setTags] = useState([]);
  useEffect(() => {
    if (!isEmpty(cronHpaData)) {
      if (tags.length) {
        const newCronHpaList = [];
        cronHpaData.records.forEach(item => {
          if (item.metadata.name.indexOf(tags[0].values[0].name) !== -1) {
            newCronHpaList.push(item);
          }
        });
        setCronHpaList(newCronHpaList);
      } else {
        setCronHpaList(cronHpaData.records);
      }
    }
  }, [cronHpaData, tags]);

  /**
   * 下载
   */
  const downloadHandle = useCallback((resourceList) => {
    const head = ['名称', '关联工作负载', '触发策略/实例个数'];
    function getTriggerStrategys(crons) {
      return crons.map((item, index) => {
        const { schedule, targetReplicas } = item;
        return `${schedule} ${targetReplicas}\r\n`;
      });
    }
    const rows = resourceList.map(cronHpa => {
      const triggerStrategyArr = getTriggerStrategys(cronHpa.spec.crons);
      const triggerStrategyStr = isEmpty(triggerStrategyArr) ? '-' : triggerStrategyArr.join('');
      return [cronHpa.metadata.name, `${cronHpa.spec.scaleTargetRef.kind}:${cronHpa.spec.scaleTargetRef.name}`, triggerStrategyStr];
    });
    downloadCsv(rows, head, 'tke_cronhpa_' + new Date().getTime() + '.csv');
  }, []);

  return (
    <Layout>
      <Body>
        <Content>
          <Content.Body>
            {
              isCronHpaInstalled === false && <TipInfo>该集群未安装CroHPA组件，请前往<a href={`/tkestack/addon?rid=1&clusterId=${clusterId}`}>扩展组件</a>进行安装</TipInfo>
            }

            <Table.ActionPanel>
              <Justify
                left={
                  <Button
                    type="primary"
                    onClick={() => {
                      delete route.queries.HPAName;
                      router.navigate({ ...urlParams, mode: 'create' }, route.queries);
                    }}
                    disabled={!isCronHpaInstalled}
                  >
                    <Trans>新建</Trans>
                  </Button>
                }
                right={(isCronHpaInstalled === true || !projectName) &&
                    (
                      <>
                        <NamespaceSelect namespaces={namespaces} />
                        <div style={{ width: 350, display: 'inline-block' }}>
                          <TagSearchBox
                            className="myTagSearchBox"
                            attributes={attributes}
                            value={tags}
                            onChange={tags => {
                              setTags(tags);
                            }}
                          />
                        </div>
                        <Button
                          onClick={() => {
                            triggerRefresh();
                          }}
                          icon="refresh"
                        />
                        <Button
                          onClick={() => {
                            downloadHandle(cronHpaList);
                          }}
                          icon="download"
                        />
                      </>
                  )
                }
              />
            </Table.ActionPanel>
            <Card>
              <Table
                records={cronHpaList}
                recordKey="id"
                columns={[
                  {
                    key: 'name',
                    header: t('名称'),
                    render: hpa => (
                      <>
                        <LinkButton onClick={() => {
                          router.navigate(
                    { ...urlParams, mode: 'detail' },
                    { ...route.queries, namespaceValue, HPAName: hpa.metadata.name });
                        }}>
                          <Text id={'hpaName' + hpa.id}>{hpa.metadata.name}</Text>
                        </LinkButton>
                        <Clip target={'#hpaName' + hpa.id} />
                      </>
                    ),
                  },
                  {
                    key: 'workload',
                    header: t('关联工作负载'),
                    render: hpa => (
                      <>
                        <Text id={'hpaWorkload' + hpa.id}>{hpa.spec.scaleTargetRef.kind}:{hpa.spec.scaleTargetRef.name}</Text>
                        <Clip target={'#hpaWorkload' + hpa.id} />
                      </>
                    ),
                  },
                  {
                    key: 'triggerStrategy',
                    header: (
                      <>
                        <Trans>触发策略/实例数</Trans>
                        <Bubble content={t('根据设置的Crontab（Crontab语法格式，例如 "0 23 * * 5"表示每周五23:00）周期性地设置实例数量')}>
                          <Icon
                            type="info"
                            style={{ marginLeft: '5px', cursor: 'pointer', verticalAlign: 'text-bottom' }}
                          />
                        </Bubble>
                      </>
                    ),
                    render: hpa => {
                      return hpa.spec.crons.map((item, index) => {
                        const { schedule, targetReplicas } = item;
                        return <Text key={index} parent="div">{`${schedule} ${targetReplicas}个`}</Text>;
                      });
                    },
                  },
                  {
                    key: 'action',
                    header: t('操作'),
                    render: (hpa) => (
                      <>
                        <LinkButton
                          tipDirection="left"
                          onClick={() => {
                    router.navigate(
                        { ...urlParams, mode: 'modify-cronhpa' },
                        { ...route.queries, namespaceValue, HPAName: hpa.metadata.name });
                          }}
                        >
                          <Trans>修改配置</Trans>
                        </LinkButton>
                        <LinkButton
                          tipDirection="left"
                          onClick={() => {
                            router.navigate(
                                { ...urlParams, mode: 'modify-yaml' },
                                { ...route.queries, namespaceValue, HPAName: hpa.metadata.name });
                          }}
                        >
                          <Trans>编辑YAML</Trans>
                        </LinkButton>
                        <LinkButton
                          tipDirection="left"
                          onClick={() => {
                            setRemoveCronHpaName(hpa.metadata.name);
                            toggle();
                          }}
                        >
                          <Trans>删除</Trans>
                        </LinkButton>
                      </>
                    )
                  }
                ]}
                addons={[
                  autotip({
                    isLoading: loading,
                    emptyText: emptyText
                  }),
                ]}
              />
              <Pagination
                recordCount={cronHpaList.length}
                pageSizeVisible={false}
              />
            </Card>
            <Modal visible={isShowing} caption="删除资源" onClose={toggle}>
              <Modal.Body><Trans>您确定要删除CronHPA: {{ removeCronHpaName }}吗？</Trans></Modal.Body>
              <Modal.Footer>
                <Button type="primary" onClick={async () => {
                  const isRemove = await deleteCronHpa({ namespace: namespaceValue, clusterId, name: removeCronHpaName });
                  if (isRemove) {
                    toggle();
                    triggerRefresh();
                  }
                }}>
                  确定
                </Button>
                <Button type="weak" onClick={toggle}>
                  取消
                </Button>
              </Modal.Footer>
            </Modal>
          </Content.Body>
        </Content>
      </Body>
    </Layout>
  );
});

export default List;
