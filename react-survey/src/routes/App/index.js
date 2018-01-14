import React, { PropTypes } from 'react';

function App({ children }) {
  return (
    <div>
      {children}
    </div>
  );
}

App.propTypes = {
  children: PropTypes.node.isRequired,
};

export default App;
