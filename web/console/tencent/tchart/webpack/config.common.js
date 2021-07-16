const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CleanWebpackPlugin = require('clean-webpack-plugin');
const Visualizer = require('webpack-visualizer-plugin');


module.exports = {
  entry: {
    TChart: './src/demo.tsx',
  },
  output: {
    path: path.resolve(__dirname, '../build'),
    filename: '[name].js'
  },
  module: {
    rules: [{
        test: /\.ts|\.tsx$/,
        loader: "awesome-typescript-loader",
      },
      {
        enforce: "pre",
        test: /\.js$/,
        loader: "source-map-loader"
      },
      {
        test: /\.less$/,
        use: [{
            loader: 'style-loader' // creates style nodes from JS strings
          },
          {
            loader: 'css-loader' // translates CSS into CommonJS
          },
          {
            loader: 'less-loader' // compiles Less to CSS
          }
        ]
      }
    ]
  },
  resolve: {
    alias: {
      core: path.resolve(__dirname, '../src/core/'),
      charts: path.resolve(__dirname, '../src/charts/'),
      "tea-components": path.resolve(__dirname, '../src/tea-components/'),
      // react 和 react-dom 控制台通过全局变量提供，我们不打包
      'react': path.resolve(__dirname, './alias/react.js'),
      'react-dom': path.resolve(__dirname, './alias/react-dom.js'),
    },
    extensions: ['.tsx', '.ts', '.js', '.json']
  },
  plugins: [
    new CleanWebpackPlugin({
      // Simulate the removal of files
      // default: false
      dry: true,

      // Write Logs to Console
      // default: false
      verbose: true,

      // Automatically remove all unused webpack assets on rebuild
      // default: true
      cleanStaleWebpackAssets: false,

      // Do not allow removal of current webpack assets
      // default: true
      protectWebpackAssets: false,

      // default: ['**/*']
      cleanOnceBeforeBuildPatterns: [path.join(process.cwd(), 'build/**/*'), '!static-files*'],
      cleanOnceBeforeBuildPatterns: [], // disables cleanOnceBeforeBuildPatterns

      // Removes files after every build (including watch mode) that match this pattern.
      // Used for files that are not created directly by Webpack.
      // Use !negative patterns to exclude files
      // default: disabled
      cleanAfterEveryBuildPatterns: ['static*.*', '!static1.js'],

      // Allow clean patterns outside of process.cwd()
      // requires dry option to be explicitly set
      // default: false
      dangerouslyAllowCleanPatternsOutsideProject: true,
    }),
    new HtmlWebpackPlugin({
      title: 'TChart',
      template: 'index.html'
    }),
    new Visualizer({
      filename: './statistics.html'
    }),
  ],
  stats: {
    colors: true,
  },
};