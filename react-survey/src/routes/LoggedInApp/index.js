import React, { PropTypes } from 'react';
import { Grid } from 'react-bootstrap';
import styles from './styles.scss';

import Sidebar from '../../components/Sidebar';
import Masthead from '../../components/Masthead';

function LoggedInApp({ children }) {
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

LoggedInApp.propTypes = {
  children: PropTypes.node.isRequired,
};

export default LoggedInApp;
