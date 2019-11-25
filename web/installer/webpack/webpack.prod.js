const path = require('path');
const webpack = require('webpack');
const HappyPack = require('happypack');
const os = require('os');
const happyThreadPool = HappyPack.ThreadPool({
  size: os.cpus().length
});

module.exports = {
  mode: 'production',
  entry: {
    app: ['./index.tsx']
  },

  output: {
    path: path.resolve(__dirname, '../public'),
    filename: 'static/js/installer.js'
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
        //exclude: [path.resolve(__dirname, "../node_modules")]
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
      id: 'happyBabel',
      loaders: ['babel-loader'],
      threadPool: happyThreadPool
    }),

    new HappyPack({
      id: 'happyCSS',
      loaders: ['style-loader', 'css-loader'],
      threadPool: happyThreadPool
    }),

    new webpack.ProgressPlugin(function handler(percentage, msg) {
      let line = 'Processing ' + msg + '\t' + Math.floor(percentage * 100) + '%';
      console.log(line);
    })
  ],

  resolve: {
    extensions: ['.tsx', '.ts', '.js', '.jsx', '.json', 'css'],
    alias: {
      '@tea/app': path.resolve(__dirname, '../node_modules/@tencent/tea-app'),
      '@tea/app/*': path.resolve(__dirname, '../node_modules/@tencent/tea-app/lib/*'),
      '@tea/component': path.resolve(__dirname, '../node_modules/@tencent/tea-component/lib'),
      '@tea/component/*': path.resolve(__dirname, '../node_modules/@tencent/tea-component/lib/*'),
      '@tencent/qcloud-lib': path.resolve(__dirname, '../node_modules/@tencent/qcloud-lib')
    }
  },

  externals: {
    react: 'window.React16',
    'react-dom': 'window.ReactDOM16'
  }
};
