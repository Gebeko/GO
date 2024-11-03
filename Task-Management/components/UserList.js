import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { getUsers, deleteUser } from './services/api';

const UserList = () => {
  const [users, setUsers] = useState([]);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const response = await getUsers();
        setUsers(response.data);
      } catch (error) {
        console.error("Error fetching users:", error);
      }
    };

    fetchUsers();
  }, []);

  const handleDelete = async (id) => {
    try {
      await deleteUser(id);
      setUsers(users.filter(user => user.id !== id));
    } catch (error) {
      console.error("Error deleting user:", error);
    }
  };

  return (
    <div>
      <h2>User List</h2>
      <ul style={{ listStyleType: 'none', padding: 0 }}>
        {users.map(user => (
          <li key={user.id} style={{ marginBottom: '10px' }}>
            {user.name} - Age: {user.age}
            <Link to={`/users/${user.id}`} style={{ marginLeft: '10px' }}>Edit</Link>
            <button onClick={() => handleDelete(user.id)} style={{ marginLeft: '10px' }}>Delete</button>
          </li>
        ))}
      </ul>
      <Link to="/add-user">Add New User</Link>
    </div>
  );
};

export default UserList;
