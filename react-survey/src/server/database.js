import Sequelize from 'sequelize';

const con = new Sequelize(process.env.DB_NAME, process.env.DB_USERNAME, process.env.DB_PASSWORD, {
  dialect: 'mysql',
  host: process.env.DB_HOST,
});

const User = con.define('user', {
  username: {
    type: Sequelize.TEXT,
    allowNull: false,
  },
  password: {
    type: Sequelize.TEXT,
    allowNull: false,
  },
});

const Study = con.define('study', {
  title: {
    type: Sequelize.TEXT,
    allowNull: false,
  },
  surveyNum: {
    type: Sequelize.INTEGER,
    allowNull: false,
  },
  partNum: {
    type: Sequelize.INTEGER,
    allowNull: false,
  },
  partCom: {
    type: Sequelize.INTEGER,
    allowNull: false,
  },
  completed: {
    type: Sequelize.INTEGER,
    allowNull: false,
  },
  daysRemaining: {
    type: Sequelize.INTEGER,
    allowNull: false,
  },
});

User.hasMany(Study);
Study.belongsTo(User);

con.sync();

export default con;
