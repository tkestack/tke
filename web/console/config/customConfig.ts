interface ICustomConfig {
  key: string;
  des: string;
  data: {
    visible: boolean;
    [key: string]: any;
  };
  children?: ICustomConfig[];
}

let CustomConfig: ICustomConfig = null;

const defaultConfig = {
  key: 'root',
  des: '整个项目的配置',
  data: {
    visible: true,
    logoDir: 'default',
    title: 'TKEStack'
  }
};

export function getCustomConfig() {
  if (CustomConfig) return CustomConfig;

  try {
    CustomConfig = JSON.parse(window?.['__CUSTOM_CONFIG']);
  } catch (error) {
    console.log('__CUSTOM_CONFIG error:', error);
    console.log('__CUSTOM_CONFIG value: ', window?.['__CUSTOM_CONFIG']);

    CustomConfig = defaultConfig;
  }

  return CustomConfig;
}
