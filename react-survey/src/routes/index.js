import React from 'react';
import Route from 'react-router/lib/Route';
import IndexRoute from 'react-router/lib/IndexRoute';
import App from './App';

// Webpack 2 supports ES2015 `import()` by auto-
// chunking assets. Check out the following for more:
// https://webpack.js.org/guides/migrating/#code-splitting-with-es2015

const importLoggedInApp = (nextState, cb) => {
  import(/* webpackChunkName: "home" */ './LoggedInApp')
    .then(module => cb(null, module.default))
    .catch(e => {
      throw e;
    });
};

const importLogin = (nextState, cb) => {
  import(/* webpackChunkName: "home" */ './Login')
    .then(module => cb(null, module.default))
    .catch(e => {
      throw e;
    });
};

const importSignup = (nextState, cb) => {
  import(/* webpackChunkName: "home" */ './Signup')
    .then(module => cb(null, module.default))
    .catch(e => {
      throw e;
    });
};

const importStudies = (nextState, cb) => {
  import(/* webpackChunkName: "home" */ './Studies')
    .then(module => cb(null, module.default))
    .catch(e => {
      throw e;
    });
};

const importTools = (nextState, cb) => {
  import(/* webpackChunkName: "tools" */ './Tools')
    .then(module => cb(null, module.default))
    .catch(e => {
      throw e;
    });
};

const importMutationTests = (nextState, cb) => {
  import(/* webpackChunkName: "home" */ './MutationTests')
    .then(module => cb(null, module.default))
    .catch(e => {
      throw e;
    });
};

// We use `getComponent` to dynamically load routes.
// https://github.com/reactjs/react-router/blob/master/docs/guides/DynamicRouting.md
const routes = (
  <Route path="/" component={App}>
    <IndexRoute getComponent={importLogin} />
    <Route path="signup" getComponent={importSignup} />
    <Route path="loggedIn" getComponent={importLoggedInApp}>
      <IndexRoute getComponent={importStudies} />
      <Route path="tools" getComponent={importTools} />
      <Route path="mutations" getComponent={importMutationTests} />
    </Route>
  </Route>
);

// Unfortunately, HMR breaks when we dynamically resolve
// routes so we need to require them here as a workaround.
// https://github.com/gaearon/react-hot-loader/issues/288
if (module.hot) {
  require('./Login'); // eslint-disable-line global-require
  require('./Signup'); // eslint-disable-line global-require
  require('./Studies'); // eslint-disable-line global-require
  require('./Tools'); // eslint-disable-line global-require
  require('./MutationTests'); // eslint-disable-line global-require
}

export default routes;
