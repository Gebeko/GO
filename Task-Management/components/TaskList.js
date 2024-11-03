import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { getTasks, deleteTask, updateTask, getUsers } from './services/api'; // Ensure updateTask is imported

const TaskList = () => {
  const [tasks, setTasks] = useState([]);
  const [selectedUserId, setSelectedUserId] = useState(null);
  const [users, setUsers] = useState([]); // Fetch users for selection

  useEffect(() => {
    const fetchTasks = async () => {
      try {
        const response = await getTasks();
        setTasks(response.data);
      } catch (error) {
        console.error("Error fetching tasks:", error);
      }
    };

    const fetchUsers = async () => {
      try {
        const response = await getUsers();
        setUsers(response.data);
      } catch (error) {
        console.error("Error fetching users:", error);
      }
    };

    fetchTasks();
    fetchUsers(); // Fetch users when the component mounts
  }, []);

  // Function to handle task deletion
  const handleDelete = async (id) => {
    try {
      await deleteTask(id);
      setTasks(tasks.filter(task => task.id !== id));
    } catch (error) {
      console.error("Error deleting task:", error);
    }
  };

  // Function to handle assigning a user to a task
  const handleAssignUser = async (taskId) => {
    if (!selectedUserId) {
      alert("Please select a user to assign.");
      return;
    }
  
    // Find the task to update
    const taskToUpdate = tasks.find(task => task.id === taskId);
  
    // Create the updated task object
    const updatedTask = {
      title: taskToUpdate.title,
      description: taskToUpdate.description,
      status: taskToUpdate.status,
      users: [...taskToUpdate.users, { id: selectedUserId }] // Add the selected user to the existing users
    };
  
    console.log("Updated Task:", updatedTask);
  
    try {
      await updateTask(taskId, updatedTask); // Call the update function
  
      // Update the local state after successful assignment
      const updatedTasks = tasks.map(task =>
        task.id === taskId ? { ...task, users: updatedTask.users } : task
      );
      setTasks(updatedTasks);
      setSelectedUserId(null); // Reset user selection
    } catch (error) {
      console.error("Error assigning user to task:", error.response.data);
    }
  };
  
  

  return (
    <div>
      <h2>Task List</h2>
      <ul style={{ listStyleType: 'none', padding: 0 }}>
        {tasks.map(task => (
          <li key={task.id} style={{ marginBottom: '10px' }}>
            <strong>{task.title}</strong> - Status: {task.status}
            <br />
            Description: {task.description}
            <br />
            Assigned User: {task.users.length > 0 ? task.users[0].name : 'None'}
            <Link to={`/tasks/${task.id}`} style={{ marginLeft: '10px' }}>Edit</Link>
            <button onClick={() => handleDelete(task.id)} style={{ marginLeft: '10px' }}>Delete</button>

            <select
              value={selectedUserId || ''}
              onChange={(e) => setSelectedUserId(e.target.value)}
              style={{ marginLeft: '10px' }}
            >
              <option value="" disabled>Select User</option>
              {users.map(user => (
                <option key={user.id} value={user.id}>{user.name}</option>
              ))}
            </select>
            <button onClick={() => handleAssignUser(task.id)} style={{ marginLeft: '10px' }}>
              Assign User
            </button>
          </li>
        ))}
      </ul>
      <Link to="/add-task">Add New Task</Link>
    </div>
  );
};

export default TaskList;
