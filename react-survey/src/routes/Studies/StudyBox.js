import React, { PropTypes } from 'react';
import { Pie } from 'react-chartjs-2';
import styles from './styles.scss';

// StudyBox component to be used in the Studies component.
// Takes in study data from database
function StudyBox({ title, surveys, participants, completed, daysRemaining }) {
  const remaining = 100 - completed;
  const pieData = {
    datasets: [
      {
        data: [completed, remaining],
        backgroundColor: ['#ea9f3c', '#fee3bf'],
      },
    ],
  };

  return (
    <div className={styles.studyBox}>
      <div className={styles.studyBoxTitle}>
        {title}
      </div>
      <table className={styles.studyBoxContent}>
        <tbody>
          <tr>
            <td>
              {surveys}
            </td>
            <td>Surveys</td>
          </tr>
          <tr>
            <td>
              {participants}
            </td>
            <td>Participants</td>
          </tr>
          <tr>
            <td>
              {completed}%
            </td>
            <td>Completed</td>
          </tr>
        </tbody>
      </table>
      <div className={styles.pieWrapper}>
        <Pie data={pieData} options={{ tooltips: { enabled: false } }} />
      </div>
      <div className={styles.daysRemaining}>
        {daysRemaining} days remaining
      </div>
    </div>
  );
}

StudyBox.propTypes = {
  title: PropTypes.string.isRequired,
  surveys: PropTypes.number.isRequired,
  participants: PropTypes.number.isRequired,
  completed: PropTypes.number.isRequired,
  daysRemaining: PropTypes.number.isRequired,
};

export default StudyBox;
