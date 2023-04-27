import defaultCustomConfig from './defaultCustomConfig.json';

interface ICustomConfigItem {
  des: string;
  visible: boolean;
  children: ICustomConfigItem;
}

interface ICustomConfig {
  title: string;
  logoDir: string;
  children: ICustomConfigItem;
}

let CustomConfig: ICustomConfig = null;

export function getCustomConfig() {
  if (CustomConfig) return CustomConfig;

  try {
    CustomConfig = JSON.parse(window?.['__CUSTOM_CONFIG']);
  } catch (error) {
    console.log('__CUSTOM_CONFIG error:', error);
    console.log('__CUSTOM_CONFIG value: ', window?.['__CUSTOM_CONFIG']);

    CustomConfig = defaultCustomConfig;
  }

  return CustomConfig;
}
