import * as React from "react";
import { Icon, Card, TagSearchBox, Select, Form, Button } from "tea-component";
import { Toolbar } from "../components/toolbar";
import { FilterTableChart } from "../components/FilterTableChart";
import { CHART_PANEL, NameValueType } from "../core";
import { CHART } from "../constants";
import { TransformField } from "../helper";

require("./Filter.less");

export interface ChartFilterPanelProps {
  width?: number;
  height?: number;
  tables: Array<CHART_PANEL.TableType>;
  groupBy: Array<NameValueType>;
}

interface ChartFilterPanelState {
  tabStatus: boolean;
  activeId: string;
  chartEndIndex: number;
  startTime: Date | string;
  endTime: Date | string;
  chartList: Array<{
    id: string;
    table: string;
    field: CHART_PANEL.MetricType;
    groupBy: Array<NameValueType>;
    conditions: Array<Array<string>>;
    tooltipLabels: (value, from) => string;
  }>;
}

/**
 * 可过滤图表
 * 每个指标单独的 维度与条件查询，有tab或list展现形式
 */
export class ChartFilterPanel extends React.Component<ChartFilterPanelProps, ChartFilterPanelState> {
  ChartPageSize = 3; // panel 为 list 状态时每次滚动加载 chart 数目
  ShowDivStyle = { flexShrink: 0, width: "100%", opacity: 1, transition: "opacity .45s" };
  HiddenDivStyle = {
    height: 0,
    overflow: "hidden",
    opacity: 0,
    flexShrink: 0,
    width: "100%",
    transition: "opacity .45s"
  };
  // 缓存每次选择不同时间段的时间粒度选项
  private periodOptions = [];
  private tabs = [];

  constructor(props) {
    super(props);
    // 默认获取一天的数据
    const endTime = new Date();
    const startTime = new Date(endTime.getTime() - 1000 * 60 * 60);
    this.periodOptions = CHART_PANEL.GeneratePeriodOptions(startTime, endTime);

    let chartList = [];
    props.tables.forEach((tableInfo, i) => {
      // 一般一个table都是相同维度
      let { groupBy = [], conditions = [] } = tableInfo;
      // this.props.groupBy 为统一的维度值，tables中每个table可以自定义维度查询值
      groupBy = this.props.groupBy.concat(groupBy);

      tableInfo.fields.forEach((field, j) => {
        const id = `${field.expr}_${i}_${j}`;
        const tooltipLabels = (value, from) => {
          return field.valueLabels ? field.valueLabels(value, from) : `${TransformField(value || 0, field.thousands, 3)}`;
        };
        this.tabs.push({ id, label: field.alias });
        chartList.push({ id, table: tableInfo.table, field, groupBy, conditions, tooltipLabels });
      });
    });
    this.state = {
      tabStatus: true,
      activeId: this.tabs.length > 0 ? this.tabs[0].id : "",
      chartEndIndex: this.ChartPageSize,
      startTime,
      endTime,
      chartList
    };
  }

  componentWillReceiveProps(nextProps) {
    if (JSON.stringify(this.props.tables) !== JSON.stringify(nextProps.tables)) {
      let chartList = [];
      nextProps.tables.forEach((tableInfo, i) => {
        // 一般一个table都是相同维度
        let { groupBy = [], conditions = [] } = tableInfo;
        groupBy = nextProps.groupBy.concat(groupBy);

        tableInfo.fields.forEach((field, j) => {
          const id = `${field.expr}_${i}_${j}`;
          const tooltipLabels = value => {
            return field.valueLabels ? field.valueLabels(value) : `${TransformField(value || 0, field.thousands, 3)}`;
          };
          this.tabs.push({ id, label: field.alias });
          chartList.push({ id, table: tableInfo.table, field, groupBy, conditions, tooltipLabels });
        });
      });
      this.setState({ chartList });
    }
  }

  // 时间选择器，在初始化时会触发
  onChangeQueryTime(startTime, endTime) {
    this.periodOptions = CHART_PANEL.GeneratePeriodOptions(startTime, endTime);
    this.setState({ startTime, endTime });
  }

  onActiveTab(tab) {
    this.setState({ activeId: tab.id });
  }

  onChangeTabStatus() {
    this.setState({ tabStatus: !this.state.tabStatus });
  }

  onScrollTabPanels(e) {
    // 图表滚动加载
    let element = e.target as any;
    // 向下滚动加载
    if (Math.floor(this.state.chartEndIndex * CHART.DefaultSize.height - element.scrollTop) <= element.clientHeight) {
      if (this.state.chartEndIndex < this.state.chartList.length) {
        this.setState({ chartEndIndex: this.state.chartEndIndex + this.ChartPageSize });
      }
    }
  }

  render() {
    const { activeId, tabStatus } = this.state;

    return (
      <div className="tea-tabs__tabpanel">
        <Toolbar
          style={{ marginLeft: 20 }}
          duration={{ from: this.state.startTime, to: this.state.endTime }}
          onChangeTime={this.onChangeQueryTime.bind(this)}
        >
          {this.props.children}
        </Toolbar>

        <Card>
          <Card.Body>
            <Tabs
              tabStatus={tabStatus}
              activeId={activeId}
              tabs={this.tabs}
              onActiveTab={this.onActiveTab.bind(this)}
              onChangeTabStatus={this.onChangeTabStatus.bind(this)}
              onScroll={this.onScrollTabPanels.bind(this)}
            >
              {this.state.chartList.map((chartInfo, index) => {
                // tab 状态时判断是否为相同 activeId，非tab状态即list显示时判断是否小于chartEndIndex的下标，小于才显示
                const isActive = tabStatus ? chartInfo.id === activeId : this.state.chartEndIndex >= index;
                const style = tabStatus
                  ? chartInfo.id === activeId
                    ? this.ShowDivStyle
                    : this.HiddenDivStyle
                  : { width: "100%" };

                return (
                  <TabPanel key={`viewgrid_${index}`} isActive={isActive} tabStatus={tabStatus} style={style}>
                    <FilterTableChart
                      chartId={chartInfo.id}
                      table={chartInfo.table}
                      startTime={this.state.startTime as any}
                      endTime={this.state.endTime as any}
                      metric={chartInfo.field}
                      tooltipLabels={chartInfo.tooltipLabels}
                      dimensions={chartInfo.groupBy}
                      conditions={chartInfo.conditions}
                      periodOptions={this.periodOptions}
                    />
                  </TabPanel>
                );
              })}
            </Tabs>
          </Card.Body>
        </Card>
      </div>
    );
  }
}

/**
 * 实现 TAB 组件，实现 tab 状态显示和 list 状态显示
 */

interface TabsProps {
  activeId: string;
  tabStatus: boolean;
  tabs: Array<{ id: string; label: string }>;
  onActiveTab: (tab) => void;
  onChangeTabStatus: () => void;
  onScroll: (e) => void;
}

interface TabsState {
  offset: number;
  scrolling: boolean;
}

class Tabs extends React.Component<TabsProps, TabsState> {
  tabBarStyle = { marginRight: 0, height: 30, lineHeight: 30 };
  scrollAreaRef = null;
  buttonRef = null;
  tabListRef = null;
  activeItemRef = null;

  constructor(props) {
    super(props);
    this.scrollAreaRef = React.createRef();
    this.buttonRef = React.createRef();
    this.tabListRef = React.createRef();
    this.activeItemRef = React.createRef();
    this.state = {
      offset: 0,
      scrolling: false
    };
    this.handleScroll = this.handleScroll.bind(this);
  }

  componentDidMount() {
    this.handleScroll();
    window.addEventListener("resize", this.handleScroll);
  }

  componentWillUnmount() {
    window.removeEventListener("resize", this.handleScroll);
  }

  handleScroll() {
    const scrolling = this.getMaxOffset() > 0;
    this.setState({ scrolling });
    // 无需滚动时重置位置
    if (!scrolling) {
      this.setState({ offset: 0 });
    } else {
      this.handleActiveItemIntoView();
    }
  }

  handleActiveItemIntoView() {
    requestAnimationFrame(() => {
      if (!this.scrollAreaRef.current || !this.activeItemRef.current) {
        return;
      }
      const scrollAreaRect = this.scrollAreaRef.current.getBoundingClientRect();
      const activeItemRect = this.activeItemRef.current.getBoundingClientRect();

      const startDelta = scrollAreaRect.left - activeItemRect.left + this.buttonRef.current.clientWidth;
      const endDelta = activeItemRect.right - scrollAreaRect.right + this.buttonRef.current.clientWidth;
      if (startDelta > 0) {
        this.setState({ offset: Math.min(0, this.state.offset + startDelta) });
      } else if (endDelta > 0) {
        this.setState({ offset: Math.max(0 - this.getMaxOffset(), this.state.offset - endDelta) });
      }
    });
  }

  getStep() {
    return this.scrollAreaRef.current.clientWidth - this.buttonRef.current.clientWidth * 2;
  }

  getMaxOffset() {
    if (!this.scrollAreaRef.current || !this.tabListRef.current) {
      return 0;
    }
    if (this.scrollAreaRef.current.clientWidth >= this.tabListRef.current.clientWidth) {
      return 0;
    }
    return (
      this.tabListRef.current.clientWidth -
      (this.scrollAreaRef.current.clientWidth - this.buttonRef.current.clientWidth * 2)
    );
  }

  handleBackward() {
    this.setState({ offset: Math.min(0, this.state.offset + this.getStep()) });
  }

  handleForward() {
    this.setState({ offset: Math.max(0 - this.getMaxOffset(), this.state.offset - this.getStep()) });
  }

  render() {
    const { activeId, tabs, tabStatus } = this.props;
    const scrollAreaStyle = tabStatus ? { marginRight: 0 } : this.tabBarStyle;

    return (
      <div className="tea-tabs">
        {/* tab bar 显示tab的标签栏*/}
        <div className="tea-tabs__tabbar">
          <div
            ref={this.scrollAreaRef}
            className={`tea-tabs__scroll-area ${this.state.scrolling ? "is-scrolling" : ""}`}
            style={{ ...scrollAreaStyle, marginRight: 30 }}
          >
            <ul
              ref={this.tabListRef}
              className="tea-tabs__tablist"
              style={{
//                transition: "transform 0.2s ease-out 0s",
                transform: `translate3d(${this.state.offset}px, 0px, 0px)`
              }}
            >
              {tabStatus &&
                tabs.map(tab => {
                  return (
                    <li
                      key={tab.id}
                      ref={tab.id === activeId ? this.activeItemRef : undefined}
                      className="tea-tabs__tabitem"
                      onClick={() => this.props.onActiveTab(tab)}
                    >
                      <a className={`tea-tabs__tab ${activeId === tab.id ? "is-active" : ""}`}>{tab.label}</a>
                    </li>
                  );
                })}
            </ul>
            {tabStatus && (
              <>
                <Button
                  ref={this.buttonRef}
                  className="tea-tabs__backward"
                  type="icon"
                  icon={"arrowleft"}
                  disabled={this.state.offset >= 0}
                  onClick={this.handleBackward.bind(this)}
                />
                <Button
                  className="tea-tabs__forward"
                  type="icon"
                  icon={"arrowright"}
                  disabled={this.state.offset <= 0 - this.getMaxOffset()}
                  onClick={this.handleForward.bind(this)}
                />
              </>
            )}
          </div>
          <div className="tea-tabs__addons">
            <span>
              <Icon
                type={tabStatus ? "viewgrid" : "viewlist"}
                style={{ cursor: "pointer" }}
                onClick={this.props.onChangeTabStatus}
              />
            </span>
          </div>
        </div>

        {/* tab 内容 */}
        <div
          className={tabStatus ? "tab-status-panel" : "no-tab-status-panel"}
          style={
            tabStatus
              ? {
                  marginLeft: `${-100 * this.props.tabs.findIndex(chartInfo => chartInfo.id === this.props.activeId)}%`
                }
              : {}
          }
          onScroll={this.props.onScroll}
        >
          {this.props.children}
        </div>
      </div>
    );
  }
}

interface TabPanelProps {
  style: React.CSSProperties;
  isActive: boolean;
  tabStatus: boolean;
}

interface TabPanelState {
  hasActive: boolean;
}

/**
 * hasActive 激活后才会渲染子项组件
 */
class TabPanel extends React.Component<TabPanelProps, TabPanelState> {
  constructor(props) {
    super(props);
    this.state = {
      hasActive: props.isActive
    };
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.isActive) {
      this.setState({ hasActive: true });
    }
  }

  render() {
    const { style = {} } = this.props;

    return (
      <div style={{ ...style }} className="tea-tabs__tabpanel">
        {this.props.tabStatus ? (
          this.state.hasActive && this.props.children
        ) : this.state.hasActive ? (
          this.props.children
        ) : (
          <div style={{ width: "100%", height: CHART.DefaultSize.height }} />
        )}
      </div>
    );
  }
}
