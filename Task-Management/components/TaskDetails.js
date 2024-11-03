import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { getTaskById } from '../../midterm/task-management/src/services/api';

const TaskDetails = () => {
  const { id } = useParams();
  const [task, setTask] = useState(null);

  useEffect(() => {
    // Fetch task details from the backend
    getTaskById(id)
      .then(response => {
        setTask(response.data);
      })
      .catch(error => {
        console.error("There was an error fetching the task!", error);
      });
  }, [id]);

  if (!task) return <div>Loading...</div>;

  return (
    <div>
      <h2>{task.title}</h2>
      <p>{task.description}</p>
    </div>
  );
};

export default TaskDetails;
