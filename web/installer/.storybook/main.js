const path = require('path');

module.exports = {
  stories: ['./stories/**/*.stories.@(js|jsx|ts|tsx|mdx)'],

  addons: ['@storybook/addon-links', '@storybook/addon-essentials'],

  webpackFinal: async (config, { configType }) => {
    // `configType` has a value of 'DEVELOPMENT' or 'PRODUCTION'
    // You can change the configuration based on that.
    // 'PRODUCTION' is used when building the static version of storybook.

    // Make whatever fine-grained changes you need
    Object.assign(config.resolve.alias, {
      '@tencent/ff-validator': path.resolve(__dirname, '../lib/ff-validator'),
      '@tencent/ff-validator/*': path.resolve(__dirname, '../lib/ff-validator/*'),
      '@tencent/ff-redux': path.resolve(__dirname, '../lib/ff-redux'),
      '@tencent/ff-redux/*': path.resolve(__dirname, '../lib/ff-redux/*'),
      '@tencent/ff-component': path.resolve(__dirname, '../lib/ff-component'),
      '@tencent/ff-component/*': path.resolve(__dirname, '../lib/ff-component/*'),
      '@tencent/tea-component': path.resolve(__dirname, '../node_modules/tea-component')
    });

    // Return the altered config
    return config;
  }
};
