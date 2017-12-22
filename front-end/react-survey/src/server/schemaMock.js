const typeDefs = `
  type Query {
    studies: [Study]
    users: [User]
  }

  # A user study.
  type Study {
    id: ID!
    title: String!
    surveyNum: Int!
    partNum: Int!
    partCom: Int!
    completed: Int!
    daysRemaining: Int!
  }

  # A user.
  type User {
    id: ID!
    username: String!
  }
`;

export default typeDefs;
