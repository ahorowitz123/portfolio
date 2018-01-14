import React, { Component, PropTypes } from 'react';
import { withApollo } from 'react-apollo';
import ApolloClient from 'apollo-client';
import gql from 'graphql-tag';
import {
  PageHeader,
  Grid,
  FormGroup,
  ControlLabel,
  FormControl,
  Button,
  HelpBlock,
} from 'react-bootstrap';

class Signup extends Component {
  state = {
    username: '',
    password1: '',
    password2: '',
    userExists: false,
  };

  handleUsernameChange = e => {
    const username = e.target.value;
    this.setState({ username });
    this.userExistsValidation(username);
  };

  handlePassword1Change = e => {
    this.setState({ password1: e.target.value });
  };

  handlePassword2Change = e => {
    this.setState({ password2: e.target.value });
  };

  userValidation() {
    const { username } = this.state;
    if (username.length < 6 || username.length > 10) {
      return false;
    }

    if (!username.match(/^[a-z0-9]+$/i)) {
      return false;
    }

    return true;
  }

  userExistsValidation(username) {
    const userQuery = gql`
      query user($username: String!) {
        user(username: $username) {
          username
        }
      }
    `;

    this.props.client
      .query({
        query: userQuery,
        variables: { username },
      })
      .then(output => {
        if (output.data.user) {
          this.setState({ userExists: true });
        } else {
          this.setState({ userExists: false });
        }
      });
  }

  password1Validation() {
    const { password1 } = this.state;
    if (password1.length < 8 || password1.length > 64) {
      return false;
    }

    if (!password1.match(/^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%^&*])/)) {
      return false;
    }

    return true;
  }

  password2Validation() {
    const { password1, password2 } = this.state;
    return password1 === password2;
  }

  getUserValidationState() {
    if (this.state.username.length === 0) {
      return null;
    }

    if (this.userValidation() && !this.state.userExists) {
      return 'success';
    }

    return 'error';
  }

  getPassword1ValidationState() {
    if (this.state.password1.length === 0) {
      return null;
    }

    if (this.password1Validation()) {
      return 'success';
    }

    return 'error';
  }

  getPassword2ValidationState() {
    if (this.state.password2.length === 0) {
      return null;
    }

    if (this.password2Validation()) {
      return 'success';
    }

    return 'error';
  }

  render() {
    const { client } = this.props;

    const formValidation = () =>
      this.userValidation() &&
      this.password1Validation() &&
      this.password2Validation() &&
      !this.state.userExists;

    const addUserMutation = gql`
      mutation addUser($username: String!, $password: String!) {
        addUser(username: $username, password: $password)
      }
    `;

    const handleOnClick = () => {
      if (formValidation()) {
        client
          .mutate({
            mutation: addUserMutation,
            variables: {
              username: this.state.username,
              password: this.state.password1,
            },
          })
          .catch(error => error);
      }
    };

    return (
      <Grid>
        <PageHeader>Sign up here!</PageHeader>
        <form>
          <FormGroup controlId="username" validationState={this.getUserValidationState()}>
            <ControlLabel>Username</ControlLabel>
            <FormControl
              type="text"
              value={this.state.username}
              onChange={this.handleUsernameChange}
            />
            {!this.userValidation() &&
              <HelpBlock>
                The username must be alpahnumeric and between 6 and 10 characters
              </HelpBlock>}
            {this.state.userExists && <HelpBlock>This username already exists</HelpBlock>}
          </FormGroup>
          <FormGroup controlId="password1" validationState={this.getPassword1ValidationState()}>
            <ControlLabel>Password</ControlLabel>
            <FormControl
              type="password"
              value={this.state.password1}
              onChange={this.handlePassword1Change}
            />
            {!this.password1Validation() &&
              <HelpBlock>
                The password must be between 8 and 64 characters and must contain at least 1
                lowercase, uppercase, number, and special character
              </HelpBlock>}
          </FormGroup>
          <FormGroup controlId="password2" validationState={this.getPassword2ValidationState()}>
            <ControlLabel>Password Again</ControlLabel>
            <FormControl
              type="password"
              value={this.state.password2}
              onChange={this.handlePassword2Change}
            />
            {!this.password2Validation() && <HelpBlock>Your passwords must match</HelpBlock>}
          </FormGroup>

          <Button type="button" onClick={handleOnClick} disabled={!formValidation()}>
            Sign up
          </Button>
        </form>
      </Grid>
    );
  }
}

Signup.propTypes = {
  client: PropTypes.instanceOf(ApolloClient).isRequired,
};

export default withApollo(Signup);
