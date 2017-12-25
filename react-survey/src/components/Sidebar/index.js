import React from 'react';
import SideNav, { Nav, NavIcon, NavText } from 'react-sidenav';

import TiClipboard from 'react-icons/lib/ti/clipboard';
import TiInputCheckedOutline from 'react-icons/lib/ti/input-checked-outline';

import styles from './styles.scss';

const Sidebar = () =>
  <div>
    <div className={styles.sidebar}>
      <SideNav highlightColor="#365963" highlightBgColor="#c7e4e8" defaultSelected="sales">
        <Nav id="studies">
          <NavIcon>
            <TiClipboard />
          </NavIcon>
          <NavText className={styles.navtext}> STUDIES </NavText>
        </Nav>
        <Nav id="surveys">
          <NavIcon>
            <TiInputCheckedOutline />
          </NavIcon>
          <NavText className={styles.navtext}> SURVEYS </NavText>
        </Nav>
      </SideNav>
    </div>
  </div>;

export default Sidebar;
