import Sequelize from 'sequelize';

const con = new Sequelize('database', 'username', 'password', {
  dialect: 'mysql',
  host: 'host',
});

const User = con.define('user', {
  username: {
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
