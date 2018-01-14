import React, { PropTypes } from 'react';
import { withApollo } from 'react-apollo';
import ApolloClient from 'apollo-client';
import gql from 'graphql-tag';
import { Link } from 'react-router';
import {
  Grid,
  Col,
  Row,
  PageHeader,
  FormGroup,
  ControlLabel,
  FormControl,
  Button,
} from 'react-bootstrap';

function Login({ client }) {
  const handleOnClick = () => {
    client
      .query({
        query: gql`
          query VerifyUser($username: String!) {
            login(username: $username)
          }
        `,
        variables: { username: this.input.value },
      })
      .then(output => {
        if (output.data.login) {
          window.location = '/LoggedIn';
        }
      });
  };

  return (
    <Grid>
      <PageHeader>Welcome to ArtsResearch!</PageHeader>
      <form>
        <FormGroup controlId="login">
          <Row>
            <Col xs={4}>
              <ControlLabel>Login</ControlLabel>
              <FormControl
                type="text"
                placeholder="Username"
                inputRef={ref => {
                  this.input = ref;
                }}
              />
            </Col>
          </Row>
        </FormGroup>
        <Row>
          <Col xs={4}>
            <Button type="button" onClick={handleOnClick}>
              Login
            </Button>
          </Col>
        </Row>
      </form>
      <Row>
        <Col xs={4}>
          <Link to="/signup">
            <Button bsStyle="link">Sign Up!</Button>
          </Link>
        </Col>
      </Row>
    </Grid>
  );
}

Login.propTypes = {
  client: PropTypes.instanceOf(ApolloClient).isRequired,
};

export default withApollo(Login);
