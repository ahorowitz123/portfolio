import React from 'react';
import { render } from 'react-dom';
import { AppContainer } from 'react-hot-loader';
import { ApolloClient, ApolloProvider, createNetworkInterface } from 'react-apollo';
import Root from './Root';

const root = document.querySelector('#root');

const networkInterface = createNetworkInterface({
  uri: 'http://localhost:3000/graphql',
});

const client = new ApolloClient({
  networkInterface,
});

const mount = RootComponent => {
  render(
    <AppContainer>
      <ApolloProvider client={client}>
        <RootComponent />
      </ApolloProvider>
    </AppContainer>,
    root
  );
};

if (module.hot) {
  module.hot.accept('./Root', () => {
    // eslint-disable-next-line global-require,import/newline-after-import
    const RootComponent = require('./Root').default;
    mount(RootComponent);
  });
}

mount(Root);
