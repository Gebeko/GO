import React, { useEffect, useState } from 'react';
import { getTasks, createTask, updateTask, deleteTask, getUsers } from '../api';
import TaskForm from './TaskForm';
import TaskDetail from './TaskDetail';
import 'bootstrap/dist/css/bootstrap.min.css';
import './TaskList.css';

const TaskList = () => {
  const [tasks, setTasks] = useState([]);
  const [users, setUsers] = useState([]);
  const [editingTask, setEditingTask] = useState(null);
  const [selectedTask, setSelectedTask] = useState(null);
  const [showDetail, setShowDetail] = useState(false);

  useEffect(() => {
    fetchTasks();
    fetchUsers();
  }, []);

  const fetchTasks = async () => {
    const taskList = await getTasks();
    setTasks(taskList);
  };

  const fetchUsers = async () => {
    const userList = await getUsers();
    setUsers(userList);
  };

  const handleTaskSave = async (task) => {
    try {
      if (editingTask) {
        await updateTask(editingTask.id, task);
      } else {
        await createTask(task);
      }
      fetchTasks();
      setEditingTask(null);
    } catch (error) {
      console.error('Error saving task:', error.response ? error.response.data : error.message);
    }
  };

  const handleEdit = (task) => setEditingTask(task);
  const handleDelete = async (id) => {
    await deleteTask(id);
    fetchTasks();
  };

  const handleShowDetail = (task) => {
    setSelectedTask(task);
    setShowDetail(true);
  };

  const handleCloseDetail = () => {
    setShowDetail(false);
    setSelectedTask(null);
  };

  return (
    <div className="container mt-4">
      <h2 className="mb-4 text-primary">Task List</h2>
      <div className="card mb-4">
        <div className="card-body">
          <TaskForm task={editingTask} users={users} onTaskSave={handleTaskSave} onCancel={() => setEditingTask(null)} />
        </div>
      </div>
      <ul className="list-group mb-4">
        {tasks.map(task => (
          <li key={task.id} className="list-group-item d-flex justify-content-between align-items-center border rounded">
            <div className="flex-grow-1 me-3">
              <strong>{task.title}</strong> - <span className="badge bg-secondary">{task.status}</span>
            </div>
            <div>
              <button className="btn btn-info btn-sm me-2" onClick={() => handleShowDetail(task)}>View Details</button>
              <button className="btn btn-warning btn-sm me-2" onClick={() => handleEdit(task)}>Edit</button>
              <button className="btn btn-danger btn-sm" onClick={() => handleDelete(task.id)}>Delete</button>
            </div>
          </li>
        ))}
      </ul>
      {showDetail && (
        <div className="mt-4">
          <TaskDetail task={selectedTask} onClose={handleCloseDetail} />
        </div>
      )}
    </div>
  );
};

export default TaskList;
