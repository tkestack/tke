export interface TencenthubNamespace {
  name?: string;
}

export interface TencenthubChart {
  app_version?: string;
  created_at?: string;
  description?: string;
  download_url?: string;
  icon?: string;
  latest_version?: string;
  name?: string;
  updated_at?: string;
}

export interface TencenthubChartVersion {
  version?: string;
  created_at?: string;
  size?: number;
  download_url?: string;
  updated_at?: string;
}

export interface TencenthubChartReadMe {
  content?: string;
}
