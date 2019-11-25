module.exports = process.env.NODE_ENV === 'procudtion' ? require('./webpack.prod.js') : require('./webpack.dev.js');
