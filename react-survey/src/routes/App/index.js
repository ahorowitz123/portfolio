import React, { PropTypes } from 'react';
import { Grid } from 'react-bootstrap';
import styles from './styles.scss';

import Sidebar from '../../components/Sidebar';
import Masthead from '../../components/Masthead';

// App component for the shell of the survey program
// Contains Masthead and Sidebar components and uses
// Bootstap grid to layout the child component
function App({ children }) {
  return (
    <div className={styles.app}>
      <Masthead />
      <Grid>
        <div className={styles.flex}>
          <Sidebar />
          <Grid>
            {children}
          </Grid>
        </div>
      </Grid>
    </div>
  );
}

App.propTypes = {
  children: PropTypes.node.isRequired,
};

export default App;
