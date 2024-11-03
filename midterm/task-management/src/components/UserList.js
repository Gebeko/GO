import React, { useEffect, useState } from 'react';
import { getUsers, createUser, updateUser, deleteUser } from '../api';
import UserForm from './UserForm';
import 'bootstrap/dist/css/bootstrap.min.css';
import './UserList.css';

const UserList = () => {
  const [users, setUsers] = useState([]);
  const [editingUser, setEditingUser] = useState(null);

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    const userList = await getUsers();
    setUsers(userList);
  };

  const handleUserSave = async (user) => {
    try {
      if (editingUser) {
        await updateUser(editingUser.id, user);
      } else {
        await createUser(user);
      }
      fetchUsers();
      setEditingUser(null);
    } catch (error) {
      console.error('Error saving user:', error.response ? error.response.data : error.message);
    }
  };

  const handleEdit = (user) => setEditingUser(user);
  const handleDelete = async (id) => {
    await deleteUser(id);
    fetchUsers();
  };

  return (
    <div className="container mt-4">
      <h2 className="mb-4 text-primary">User List</h2>
      <div className="card mb-4">
        <div className="card-body">
          <UserForm user={editingUser} onUserSave={handleUserSave} onCancel={() => setEditingUser(null)} />
        </div>
      </div>
      <ul className="list-group mt-3">
        {users.map(user => (
          <li key={user.id} className="list-group-item d-flex justify-content-between align-items-center border rounded">
            <div>
              {user.name} ({user.age})
            </div>
            <div>
              <button className="btn btn-warning btn-sm me-2" onClick={() => handleEdit(user)}>Edit</button>
              <button className="btn btn-danger btn-sm" onClick={() => handleDelete(user.id)}>Delete</button>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default UserList;
