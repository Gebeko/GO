import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { createUser, updateUser, getUser } from './services/api';

const UserForm = () => {
  const [name, setName] = useState('');
  const [age, setAge] = useState('');
  const navigate = useNavigate();
  const { id } = useParams();

  useEffect(() => {
    if (id) {
      const fetchUser = async () => {
        const response = await getUser(id);
        setName(response.data.name);
        setAge(response.data.age);
      };
      fetchUser();
    }
  }, [id]);

  const handleSubmit = async (e) => {
    e.preventDefault();

    const userData = {
      name,
      age: parseInt(age), // Ensure age is a number
    };

    try {
      if (id) {
        await updateUser(id, userData);
      } else {
        await createUser(userData);
      }
      navigate('/'); // Redirect to the dashboard after submission
    } catch (error) {
      console.error("Error adding/updating user:", error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h2>{id ? 'Edit User' : 'Add User'}</h2>
      <div>
        <label>Name:</label>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
      </div>
      <div>
        <label>Age:</label>
        <input
          type="number"
          value={age}
          onChange={(e) => setAge(e.target.value)}
          required
        />
      </div>
      <button type="submit">{id ? 'Update User' : 'Add User'}</button>
    </form>
  );
};

export default UserForm;
