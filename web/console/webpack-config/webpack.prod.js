const baseConfig = require('./webpack.base');

module.exports = ({ version }) => baseConfig({ version, mode: 'production' });
