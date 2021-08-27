/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */
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
