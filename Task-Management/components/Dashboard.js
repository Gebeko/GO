import React from 'react';
import TaskList from './TaskList';
import UserList from './UserList';

const Dashboard = () => {
  return (
    <div>
      <h1>Dashboard</h1>
      <h2>Tasks</h2>
      <TaskList />
      <h2>Users</h2>
      <UserList />
    </div>
  );
};

export default Dashboard;
