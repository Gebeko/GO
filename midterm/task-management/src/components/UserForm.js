import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import 'bootstrap/dist/css/bootstrap.min.css';

const UserForm = ({ user, onUserSave, onCancel }) => {
  const [name, setName] = useState('');
  const [age, setAge] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    if (user) {
      setName(user.name);
      setAge(user.age);
    }
  }, [user]);

  const handleSubmit = (e) => {
    e.preventDefault();
    const newUser = { name, age: parseInt(age) };
    onUserSave(newUser);
    setName('');
    setAge('');
  };

  return (
    <form onSubmit={handleSubmit} className="mb-4">
      <div className="mb-3">
        <label htmlFor="name" className="form-label">Name</label>
        <input
          id="name"
          type="text"
          placeholder="Enter name"
          className="form-control"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
      </div>
      <div className="mb-3">
        <label htmlFor="age" className="form-label">Age</label>
        <input
          id="age"
          type="number"
          placeholder="Enter age"
          className="form-control"
          value={age}
          onChange={(e) => setAge(e.target.value)}
          required
        />
      </div>
      <div className="d-flex justify-content-between">
        <button type="submit" className="btn btn-primary me-2">
          {user ? 'Update User' : 'Create User'}
        </button>
        <button type="button" className="btn btn-secondary me-2" onClick={onCancel}>Cancel</button>
        <button type="button" className="btn btn-light" onClick={() => navigate('/')}>
          Back to Main Menu
        </button>
      </div>
    </form>
  );
};

export default UserForm;
