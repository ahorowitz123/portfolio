import React, { PropTypes } from 'react';
import { PageHeader, FormGroup, ControlLabel, FormControl, Button } from 'react-bootstrap';
import { withApollo } from 'react-apollo';
import ApolloClient from 'apollo-client';
import gql from 'graphql-tag';

// Login component for the login page
function Login({ client }) {
  // function to verify user information on click
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
          window.location = '/studies';
        }
      });
  };

  return (
    <div>
      <PageHeader>Welcome to ArtsResearch!</PageHeader>
      <form>
        <FormGroup controlId="login">
          <ControlLabel>Add User</ControlLabel>
          <FormControl
            type="text"
            placeholder="Username"
            inputRef={ref => {
              this.input = ref;
            }}
          />
        </FormGroup>

        <Button type="button" onClick={handleOnClick}>
          Login
        </Button>
      </form>
    </div>
  );
}

Login.propTypes = {
  client: PropTypes.instanceOf(ApolloClient).isRequired,
};

export default withApollo(Login);
