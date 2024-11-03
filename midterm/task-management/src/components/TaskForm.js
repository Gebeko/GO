import React, { useState, useEffect } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';
import { useNavigate } from 'react-router-dom';

const TaskForm = ({ task, users, onTaskSave, onCancel }) => {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [status, setStatus] = useState('Pending');
  const [assignedUsers, setAssignedUsers] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    if (task) {
      setTitle(task.title);
      setDescription(task.description);
      setStatus(task.status);
      setAssignedUsers(task.Users.map(user => user.id));
    } else {
      // Reset form state when there's no task
      resetForm();
    }
  }, [task]);

  const resetForm = () => {
    setTitle('');
    setDescription('');
    setStatus('Pending');
    setAssignedUsers([]);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const newTask = {
      title,
      description,
      status,
      Users: assignedUsers.map(id => ({ id }))
    };

    onTaskSave(newTask);
    resetForm(); // Reset the form after saving
  };

  const handleUserChange = (userId) => {
    setAssignedUsers(prev => {
      if (prev.includes(userId)) {
        return prev.filter(id => id !== userId);
      }
      return [...prev, userId];
    });
  };

  return (
    <form onSubmit={handleSubmit} className="mb-4">
      <div className="mb-3">
        <label className="form-label" htmlFor="taskTitle">Title</label>
        <input
          id="taskTitle"
          type="text"
          className="form-control"
          placeholder="Title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
        />
      </div>
      <div className="mb-3">
        <label className="form-label" htmlFor="taskDescription">Description</label>
        <textarea
          id="taskDescription"
          className="form-control"
          placeholder="Description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          required
        />
      </div>
      <div className="mb-3">
        <label className="form-label" htmlFor="taskStatus">Status</label>
        <select
          id="taskStatus"
          className="form-select"
          value={status}
          onChange={(e) => setStatus(e.target.value)}
        >
          <option value="Pending">Pending</option>
          <option value="In Progress">In Progress</option>
          <option value="Completed">Completed</option>
        </select>
      </div>
      <div className="mb-3">
        <h4>Assign Users:</h4>
        {users.map(user => (
          <div key={user.id} className="form-check">
            <input
              className="form-check-input"
              type="checkbox"
              id={`user-${user.id}`}
              checked={assignedUsers.includes(user.id)}
              onChange={() => handleUserChange(user.id)}
            />
            <label className="form-check-label" htmlFor={`user-${user.id}`}>
              {user.name}
            </label>
          </div>
        ))}
      </div>
      <div className="d-flex justify-content-between">
        <button type="submit" className="btn btn-primary me-2">
          {task ? 'Update Task' : 'Create Task'}
        </button>
        <button type="button" className="btn btn-secondary me-2" onClick={() => { onCancel(); resetForm(); }}>
          Cancel
        </button>
        <button type="button" className="btn btn-light" onClick={() => navigate('/')}>
          Back to Main Menu
        </button>
      </div>
    </form>
  );
};

export default TaskForm;
