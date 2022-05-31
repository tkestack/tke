import React, { useEffect, useState } from 'react';
import { Modal, Segment, DatePicker, Text, Select, Row, Col, Card, Table, H3, StatusTip } from 'tea-component';
import { BasicLine } from 'tea-chart';
import { fetchMetrics } from '@src/webApi/monitor';
import { useFetch } from '@src/modules/common/hooks';
import {
  IMonitorPanelProps,
  dateRangeTypeOptions,
  timeGranularityOptions,
  DateRangeTypeEnum,
  IChartRenderProps
} from './constants';
import moment from 'moment';
import { isNumber } from 'lodash';

const { RangePicker } = DatePicker;
const { selectable, scrollable } = Table.addons;

export const MonitorPanel = ({
  title,
  conditions,
  groups,
  instanceType,
  visible,
  onClose,
  instanceList,
  defaultSelectedInstances
}: IMonitorPanelProps) => {
  const [dateRangeType, setDataRangeType] = useState(DateRangeTypeEnum.HOUR);
  const [dateRange, setDataRange] = useState({
    startTime: moment().subtract(1, 'hour'),
    endTime: moment()
  });

  // 统计粒度
  const [timeGran, setTimeGran] = useState(`timestamp(${1 * 60}s)`);

  const [selectedInstances, setSelectedInstances] = useState([]);

  useEffect(() => {
    setSelectedInstances(defaultSelectedInstances);
  }, [defaultSelectedInstances]);

  function handleDateRangeChange(rangeType: DateRangeTypeEnum, range?: [moment.Moment, moment.Moment]) {
    let startTime = moment();
    let endTime = moment();
    switch (rangeType) {
      case DateRangeTypeEnum.HOUR:
        startTime = moment().subtract(1, 'hour');
        setTimeGran(`timestamp(${1 * 60}s)`);
        break;
      case DateRangeTypeEnum.DAY:
        startTime = moment().subtract(1, 'day');
        setTimeGran(`timestamp(${1 * 60}s)`);
        break;
      case DateRangeTypeEnum.WEEK:
        startTime = moment().subtract(1, 'week');
        setTimeGran(`timestamp(${60 * 60}s)`);
        break;
      case DateRangeTypeEnum.CUSTOM:
        startTime = range[0];
        endTime = range[1];
        setTimeGran(`timestamp(${60 * 60}s)`);
        break;
    }

    setDataRangeType(rangeType);
    setDataRange({
      startTime,
      endTime
    });
  }

  return (
    <Modal caption={title} visible={visible} size={1050} onClose={onClose}>
      <Modal.Body>
        <Row>
          <Col>
            <Segment
              options={dateRangeTypeOptions}
              value={dateRangeType}
              onChange={(type: DateRangeTypeEnum) => handleDateRangeChange(type)}
            />

            <RangePicker
              showTime
              placeholder="选择日期"
              range={[moment().subtract(1, 'months'), moment()]}
              onChange={range => handleDateRangeChange(DateRangeTypeEnum.CUSTOM, range)}
            />

            <Text style={{ marginRight: 10, marginLeft: 20 }}>统计粒度:</Text>
            <Select
              appearance="button"
              matchButtonWidth
              options={timeGranularityOptions[dateRangeType]}
              value={timeGran}
              onChange={value => setTimeGran(value)}
            />
          </Col>
        </Row>

        <Row>
          <Col span={18}>
            <Card bordered>
              <Card.Body style={{ height: 500, overflow: 'auto' }}>
                {groups.map((group, index) => (
                  <ChartRender
                    key={index}
                    conditions={conditions}
                    group={group}
                    timeGran={timeGran}
                    instanceType={instanceType}
                    dateRange={dateRange}
                    selectedInstances={selectedInstances}
                  />
                ))}
              </Card.Body>
            </Card>
          </Col>

          <Col span={6}>
            <Table
              bordered
              columns={[
                {
                  key: 'name',
                  header: '名称'
                }
              ]}
              records={instanceList.map(name => ({ name }))}
              recordKey="name"
              addons={[
                selectable({
                  value: selectedInstances,
                  onChange: names => setSelectedInstances(names),
                  rowSelect: true
                }),
                scrollable({ maxHeight: 460, minHeight: 460 })
              ]}
            />
          </Col>
        </Row>
      </Modal.Body>
    </Modal>
  );
};

function ChartRender({
  conditions,
  group: { fields, by },
  timeGran,
  instanceType,
  dateRange,
  selectedInstances
}: IChartRenderProps) {
  const { data: metrics, status } = useFetch(async () => {
    const rsp = await fetchMetrics({
      conditions,
      fields: fields.map(({ expr }) => expr),
      groupBy: [timeGran, ...by],
      table: instanceType,
      startTime: dateRange.startTime.valueOf(),
      endTime: dateRange.endTime.valueOf()
    });

    return {
      data: rsp
    };
  }, [conditions, fields, by, instanceType, dateRange, timeGran]);

  const chartsDataSource = React.useMemo(() => {
    return (
      fields?.map(({ alias, unit }, cIndex) => {
        return {
          title: `${alias} (${unit})`,
          data:
            metrics?.data
              ?.filter(items => {
                const byLength = by.length;
                const nameIndex = by.findIndex(i => i === 'name');
                const byItems = items.slice(-byLength);
                const name = byItems[nameIndex];

                return selectedInstances.includes(name);
              })
              ?.map(([time, ...others]) => {
                const byLength = by.length;
                const byItems = others.slice(-byLength);
                const mark = byItems.join('-');

                const _value = others?.[cIndex];
                const value = isNumber(_value) ? +_value.toFixed(6) : _value;

                return {
                  time: moment(time).format('MM-DD HH:mm'),
                  value,
                  mark
                };
              }) ?? []
        };
      }) ?? []
    );
  }, [metrics, by, selectedInstances, fields]);

  return (
    <>
      {chartsDataSource.map(({ title, data }) => (
        <React.Fragment key={title}>
          <H3>{title}</H3>
          <BasicLine
            canvasMode
            height={250}
            size={1}
            position="time*value"
            dataSource={data}
            color="mark"
            tips={status === 'loading' ? <StatusTip status="loading" /> : undefined}
          />
        </React.Fragment>
      ))}
    </>
  );
}
