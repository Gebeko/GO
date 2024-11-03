import axios from 'axios';

const API_URL = 'http://localhost:8080';

// Task API
export const getTasks = () => axios.get(`${API_URL}/tasks`);
export const getTask = (id) => axios.get(`${API_URL}/tasks/${id}`);
export const createTask = (task) => axios.post(`${API_URL}/tasks`, task);
export const updateTask = (id, task) => axios.put(`${API_URL}/tasks/${id}`, task);
export const deleteTask = (id) => axios.delete(`${API_URL}/tasks/${id}`);

// User API
export const getUsers = () => axios.get(`${API_URL}/users`);
export const getUser = (id) => axios.get(`${API_URL}/users/${id}`);
export const createUser = (user) => axios.post(`${API_URL}/users`, user);
export const updateUser = (id, user) => axios.put(`${API_URL}/users/${id}`, user);
export const deleteUser = (id) => axios.delete(`${API_URL}/users/${id}`);
