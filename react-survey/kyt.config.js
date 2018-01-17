// Base kyt config.
// Edit these properties to make changes.

const Dotenv = require('dotenv-webpack');

module.exports = {
  reactHotLoader: true,
  debug: false,

  modifyWebpackConfig(kytConfig) {
    const appConfig = Object.assign({}, kytConfig);
    appConfig.plugins.push(new Dotenv());
    return appConfig;
  },
};
