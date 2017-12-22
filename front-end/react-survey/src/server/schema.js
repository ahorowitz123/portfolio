import {
  GraphQLInt,
  GraphQLObjectType,
  GraphQLList,
  GraphQLNonNull,
  GraphQLSchema,
  GraphQLString,
  GraphQLID,
  GraphQLBoolean,
} from 'graphql';

import Db from './database';

/* eslint-disable no-use-before-define */

let User;

const studyType = new GraphQLObjectType({
  name: 'Study',
  description: 'A user study.',
  fields: {
    id: {
      type: new GraphQLNonNull(GraphQLID),
    },
    title: {
      type: new GraphQLNonNull(GraphQLString),
    },
    surveyNum: {
      type: new GraphQLNonNull(GraphQLInt),
    },
    partNum: {
      type: new GraphQLNonNull(GraphQLInt),
    },
    partCom: {
      type: new GraphQLNonNull(GraphQLInt),
    },
    completed: {
      type: new GraphQLNonNull(GraphQLInt),
    },
    daysRemaining: {
      type: new GraphQLNonNull(GraphQLInt),
    },
  },
});

const userType = new GraphQLObjectType({
  name: 'User',
  description: 'A user.',
  fields: {
    id: {
      type: new GraphQLNonNull(GraphQLID),
    },
    username: {
      type: new GraphQLNonNull(GraphQLString),
    },
  },
});

const queryType = new GraphQLObjectType({
  name: 'Query',
  fields: {
    studies: {
      type: new GraphQLList(studyType),
      resolve: () => Db.models.study.findAll(),
    },
    users: {
      type: new GraphQLList(userType),
      resolve: () => Db.models.user.findAll(),
    },
    currentUser: {
      type: userType,
      resolve: () => User,
    },
    login: {
      type: new GraphQLNonNull(GraphQLBoolean),
      args: {
        username: {
          type: new GraphQLNonNull(GraphQLString),
        },
      },
      resolve: (root, { username }) =>
        Db.models.user.findOne({ where: { username } }).then(user => {
          if (user) {
            User = user;
            return true;
          }
          return false;
        }),
    },
  },
});

const addUser = {
  name: 'AddUser',
  description: 'Add a new user.',
  type: userType,
  args: {
    username: { type: new GraphQLNonNull(GraphQLString) },
  },
  resolve: (value, { username }) => Db.models.user.create({ username }).then(user => user),
};

const mutationType = new GraphQLObjectType({
  name: 'Mutation',
  fields: {
    addUser,
  },
});

const Schema = new GraphQLSchema({
  query: queryType,
  mutation: mutationType,
});

export default Schema;
