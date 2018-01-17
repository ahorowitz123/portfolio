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

import bcrypt from 'bcrypt';
import jwt from 'jsonwebtoken';

import Db from './database';

/* eslint-disable no-use-before-define */

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
    jwt: {
      type: new GraphQLNonNull(GraphQLString),
    },
  },
});

const queryType = new GraphQLObjectType({
  name: 'Query',
  fields: {
    studies: {
      type: new GraphQLList(studyType),
      resolve: (root, args, ctx) =>
        ctx.user.then(user => {
          if (!user) {
            return Promise.reject('Unauthorized');
          }
          return Db.models.study.findAll();
        }),
    },
    users: {
      type: new GraphQLList(userType),
      resolve: () => Db.models.user.findAll(),
    },
    user: {
      type: userType,
      args: {
        username: {
          type: new GraphQLNonNull(GraphQLString),
        },
      },
      resolve: (root, { username }) =>
        Db.models.user.findOne({ where: { username } }).then(user => user),
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
  type: new GraphQLNonNull(userType),
  args: {
    username: { type: new GraphQLNonNull(GraphQLString) },
    password: { type: new GraphQLNonNull(GraphQLString) },
  },
  resolve: (value, { username, password }, ctx) =>
    Db.models.user.findOne({ where: { username } }).then(existing => {
      if (!existing) {
        const saltRounds = 10;
        const hash = bcrypt.hashSync(password, saltRounds);
        return Db.models.user.create({ username, password: hash }).then(user => {
          user.jwt = jwt.sign({ id: user.id, username: user.username }, process.env.SECRET_KEY);
          ctx.user = Promise.resolve(user);
          return user;
        });
      }

      return Promise.reject('username already exists');
    }),
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
