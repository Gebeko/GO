import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { getTask, createTask, updateTask } from './services/api';

const TaskForm = () => {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [status, setStatus] = useState('Pending');
  const navigate = useNavigate();
  const { id } = useParams(); 

  useEffect(() => {
    if (id) {
      const fetchTask = async () => {
        try {
          const response = await getTask(id);
          setTitle(response.data.title);
          setDescription(response.data.description);
          setStatus(response.data.status);
        } catch (error) {
          console.error("Error fetching task:", error);
        }
      };

      fetchTask();
    }
  }, [id]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      if (id) {

        await updateTask(id, { title, description, status });
      } else {

        await createTask({ title, description, status });
      }
      navigate('/');
  
    } catch (error) {
      console.error("Error saving task:", error);
    }
  };

  return (
    <div>
      <h2>{id ? 'Edit Task' : 'Add New Task'}</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Title:</label>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
          />
        </div>
        <div>
          <label>Description:</label>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            required
          />
        </div>
        <div>
          <label>Status:</label>
          <select value={status} onChange={(e) => setStatus(e.target.value)}>
            <option value="Pending">Pending</option>
            <option value="In Progress">In Progress</option>
            <option value="Completed">Completed</option>
          </select>
        </div>
        <button type="submit">{id ? 'Save Changes' : 'Add Task'}</button>
      </form>
    </div>
  );
};

export default TaskForm;
