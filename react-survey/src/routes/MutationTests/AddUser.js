import React, { PropTypes } from 'react';
import { FormGroup, ControlLabel, FormControl, Button } from 'react-bootstrap';
import { graphql } from 'react-apollo';
import gql from 'graphql-tag';

function AddUser({ mutate }) {
  const handleOnClick = () => {
    mutate({
      variables: { username: this.input.value },
    })
      .then(({ data }) => data)
      .catch(error => error);
  };

  return (
    <form>
      <FormGroup controlId="addUser">
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
        Add User
      </Button>
    </form>
  );
}

AddUser.propTypes = {
  mutate: PropTypes.func.isRequired,
};

const addUser = gql`
  mutation addUser($username: String!) {
    addUser(username: $username) {
      id
    }
  }
`;

const AddUserWithData = graphql(addUser)(AddUser);

export default AddUserWithData;
