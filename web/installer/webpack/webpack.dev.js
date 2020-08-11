const path = require('path');
const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const HappyPack = require('happypack');
const ProgressBarPlugin = require('progress-bar-webpack-plugin');

const os = require('os');
const happyThreadPool = HappyPack.ThreadPool({
  size: os.cpus().length
});

module.exports = {
  devtool: 'inline-source-map',
  mode: 'development',

  entry: {
    app: ['./index.tsx']
    // vendor: ['react', 'react-dom']
  },

  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'js/[name].js',
    publicPath: 'http://localhost:6688/'
  },

  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: [
          'happypack/loader?id=happyBabel',
          {
            loader: 'ts-loader',
            options: {
              transpileOnly: true
            }
          }
        ]
        //exclude: [path.resolve(__dirname, '../node_modules')]
      },
      {
        test: /\.jsx?$/,
        use: ['happypack/loader?id=happyBabel']
      },
      {
        test: /\.css?$/,
        use: ['happypack/loader?id=happyCSS']
      }
    ]
  },

  plugins: [
    new HappyPack({
      id: 'happyTs',
      loaders: [
        {
          loader: 'ts-loader',
          options: {
            happyPackMode: true,
            transpileOnly: true
          }
        }
      ],
      threadPool: happyThreadPool
    }),

    new HappyPack({
      id: 'happyBabel',
      loaders: [
        {
          loader: 'babel-loader'
        }
      ],
      threadPool: happyThreadPool
    }),

    new HappyPack({
      id: 'happyCSS',
      loaders: ['style-loader', 'css-loader'],
      threadPool: happyThreadPool
    }),

    new ProgressBarPlugin({
      summary: false
    }),

    new HtmlWebpackPlugin({
      template: path.resolve(__dirname, '../index.html')
    }),
    new webpack.HotModuleReplacementPlugin()
  ],

  resolve: {
    extensions: ['.tsx', '.ts', '.js', '.jsx', '.json', 'css'],
    alias: {
      '@tea/app': path.resolve(__dirname, '../node_modules/@tencent/tea-app'),
      '@tea/app/*': path.resolve(__dirname, '../node_modules/@tencent/tea-app/lib/*'),
      '@tea/component': path.resolve(__dirname, '../node_modules/@tencent/tea-component/lib'),
      '@tea/component/*': path.resolve(__dirname, '../node_modules/@tencent/tea-component/lib/*'),
      '@tencent/ff-validator': path.resolve(__dirname, '../lib/ff-validator'),
      '@tencent/ff-validator/*': path.resolve(__dirname, '../lib/ff-validator/*'),
      '@tencent/ff-redux': path.resolve(__dirname, '../lib/ff-redux'),
      '@tencent/ff-redux/*': path.resolve(__dirname, '../lib/ff-redux/*'),
      '@tencent/ff-component': path.resolve(__dirname, '../lib/ff-component'),
      '@tencent/ff-component/*': path.resolve(__dirname, '../lib/ff-component/*')
    }
  },

  externals: {
    react: 'window.React16',
    'react-dom': 'window.ReactDOM16'
  }
};
