import React, { PropTypes } from 'react';
import { Row, Col } from 'react-bootstrap';
import { graphql } from 'react-apollo';
import gql from 'graphql-tag';

import StudyBox from './StudyBox';
import styles from './styles.scss';

// Component for studies page. Uses the StudyBox component
function Studies({ data: { loading, error, studies } }) {
  if (loading) {
    return <p>Loading ...</p>;
  }
  if (error) {
    return (
      <p>
        {error.message}
      </p>
    );
  }

  return (
    <div className={styles.studyBoxesWrapper}>
      <Row>
        {studies.map(study =>
          <Col className={styles.studyBoxWrapper} lg={4}>
            <StudyBox {...study} />
          </Col>
        )}
      </Row>
    </div>
  );
}

Studies.propTypes = {
  data: PropTypes.shape({
    loading: PropTypes.bool.isRequired,
    studies: PropTypes.Array,
  }).isRequired,
};

const GetStudies = gql`
  query GetStudies {
    studies {
      id
      title
      surveyNum
      partNum
      partCom
      completed
      daysRemaining
    }
  }
`;

const StudiesWithData = graphql(GetStudies)(Studies);

export default StudiesWithData;
