import { createCSRFHeader } from '@helper'

export async function apiRequest({ data }) {
  let res;
  // 融合版的tke 监控，使用influxdb
  let apiInfo = window["modules"]["monitor"];
  let json = {
    apiVersion: apiInfo
      ? `${apiInfo.groupName}/${apiInfo.version}`
      : "monitor.tkestack.io/v1",
    kind: "Metric",
    query: data.data.RequestBody,
  };
  json.query.offset = 0;
  json.query.conditions = json.query.conditions.map((item) => ({
    key: item[0],
    expr: item[1],
    value: item[2],
  }));
  let projectHeader = {};
  const projectName = new URLSearchParams(window.location.search).get(
    "projectName"
  );

  if (projectName) {
    projectHeader = {
      "X-TKE-ProjectName": projectName,
    };
  }

  try {
    res = await fetch(`/apis/monitor.tkestack.io/v1/metrics/`, {
      method: "POST",
      mode: "cors",
      cache: "no-cache",
      credentials: "same-origin",
      headers: {
        "Content-Type": "application/json",
        ...createCSRFHeader(),
        ...projectHeader,
      },
      redirect: "follow",
      referrer: "no-referrer",
      body: JSON.stringify(json),
    });
    res = await res.json();
    return JSON.parse(res.jsonResult);
  } catch (error) {
    return {
      columns: 0,
      data: [],
    };
  }
}

export async function request(data) {
  try {
    // 发送 API 请求
    let res = await apiRequest({
      data,
    });

    return { columns: res.columns, data: res.data || [] };
  } catch (error) {
    console.error(error, data);
    throw error;
  }
}

const getCookie = (name: string) => {
  let reg = new RegExp("(?:^|;+|\\s+)" + name + "=([^;]*)"),
    match = document.cookie.match(reg);

  return !match ? "" : match[1];
};
