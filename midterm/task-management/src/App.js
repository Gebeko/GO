import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import UserList from './components/UserList';
import TaskList from './components/TaskList';
import MainMenu from './components/MainMenu';

const App = () => {
  return (
    <Router>
      <div>
        <Routes>
          <Route path="/" element={<MainMenu />} />
          <Route path="/users" element={<UserList />} />
          <Route path="/tasks" element={<TaskList />} />
        </Routes>
      </div>
    </Router>
  );
};

export default App;