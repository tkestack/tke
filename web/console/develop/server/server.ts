import * as path from 'path';
import * as express from 'express';
import * as webpack from 'webpack';

export function serve(serverPort?: number) {
  let port = serverPort || process.env.PORT;

  let app = express();

  setupWebpackDevelopServer(app);

  app.use('/static', express.static('dist'));

  // app.use('/static/*', express.static(path.resolve(__dirname, '../../dist/js/*')));
  app.use('*', express.static(path.resolve(__dirname, '../../index.html')));

  app.listen(port, (err: Error) => {
    if (err) {
      console.error(JSON.stringify(err));
      return;
    }

    console.log(`\nTencent Kubernetes Engine server served at http://localhost:${port}\n\n`);
  });

  return app;
}

function setupWebpackDevelopServer(app: express.Express) {
  let config = require('../../webpack/webpack.config.js');
  let compiler = webpack(config);

  let devMiddleware = require('webpack-dev-middleware')(compiler, {
    publicPath: config.output.publicPath,
    noInfo: true,
    stats: { colors: true },
    poll: true,
    quiet: false,
    reload: true,
    writeToDisk: true
  });

  let hotMiddleware = require('webpack-hot-middleware')(compiler, { reload: true });

  app.use(devMiddleware);
  app.use(hotMiddleware);
}
