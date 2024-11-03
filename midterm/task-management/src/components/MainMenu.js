import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { getTasks } from '../api';
import 'bootstrap/dist/css/bootstrap.min.css';
import { FaUserPlus, FaTasks } from 'react-icons/fa'; // Import icons

const MainMenu = () => {
  const [tasks, setTasks] = useState([]);

  useEffect(() => {
    fetchTasks();
  }, []);

  const fetchTasks = async () => {
    try {
      const taskList = await getTasks();
      setTasks(taskList);
    } catch (error) {
      console.error('Error fetching tasks:', error);
    }
  };

  const incompleteTasks = tasks.filter(task => task.status !== 'Completed');

  return (
    <div className="container text-center my-5" style={{ backgroundColor: '#f8f9fa', borderRadius: '10px', padding: '30px' }}>
      <h1 className="mb-4 text-primary">Welcome to the Task Management App created by Baatar</h1>
      <div className="d-flex justify-content-center mb-4">
        <Link to="/users" className="me-2">
          <button className="btn btn-primary btn-lg">
            <FaUserPlus className="me-2" /> Create New User
          </button>
        </Link>
        <Link to="/tasks">
          <button className="btn btn-secondary btn-lg">
            <FaTasks className="me-2" /> Create New Task
          </button>
        </Link>
      </div>

      <h2 className="mt-5 text-secondary">Available Incomplete Tasks:</h2>
      <div className="card mt-3">
        <div className="card-body">
          <ul className="list-group">
            {incompleteTasks.length > 0 ? (
              incompleteTasks.map(task => (
                <li key={task.id} className="list-group-item">
                  <strong>{task.title}</strong> - [{task.status}]
                </li>
              ))
            ) : (
              <li className="list-group-item">No incomplete tasks available.</li>
            )}
          </ul>
        </div>
      </div>
    </div>
  );
};

export default MainMenu;
