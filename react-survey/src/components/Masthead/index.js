import React from 'react';
import Link from 'react-router/lib/Link';
import { Navbar, Nav, NavItem } from 'react-bootstrap';
import styles from './styles.scss';

function Masthead() {
  return (
    <Navbar collapseOnSelect className={styles.navbar}>
      <Navbar.Header>
        <Navbar.Brand>
          <Link to="/">ArtsResearch</Link>
        </Navbar.Brand>
        <Navbar.Toggle />
      </Navbar.Header>
      <Navbar.Collapse>
        <Nav>
          <NavItem eventKey={1} href="/">
            Home
          </NavItem>
        </Nav>
      </Navbar.Collapse>
    </Navbar>
  );
}

export default Masthead;
