import React, { createContext, useContext, useState, useEffect } from 'react';
import axiosInstance from '../Login/axiosInstance';
import {jwtDecode} from 'jwt-decode';

const ItemContext = createContext();

export const useItemContext = () => useContext(ItemContext);

export const ItemProvider = ({ children }) => {
    const [items, setItems] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [userRole, setUserRole] = useState('');

    const fetchItems = async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await axiosInstance.get('http://localhost:8080/items');
            setItems(response.data);
        } catch (error) {
            setError("Failed to fetch items. Please try again.");
            console.error("Error fetching items:", error);
        } finally {
            setLoading(false);
        }
    };

    const addItem = async (newItem) => {
        try {
            const response = await axiosInstance.post('http://localhost:8080/items', newItem);
            setItems([...items, response.data]);
        } catch (error) {
            console.error("Error adding item:", error);
        }
    };

    const deleteItem = async (id) => {
        try {
            await axiosInstance.delete(`http://localhost:8080/items/${id}`);
            setItems(items.filter(item => item.id !== id));
        } catch (error) {
            console.error("Error deleting item:", error);
        }
    };

    const updateItem = async (updatedItem) => {
        try {
            const response = await axiosInstance.put(`http://localhost:8080/items/${updatedItem.id}`, updatedItem);
            setItems(items.map(item => (item.id === updatedItem.id ? response.data : item)));
        } catch (error) {
            console.error("Error updating item:", error);
        }
    };

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (token) {
            setIsAuthenticated(true);
            const decodedToken = jwtDecode(token);
            setUserRole(decodedToken.role);
        }
        fetchItems();
    }, [isAuthenticated]);

    return (
        <ItemContext.Provider value={{ items, loading, error, addItem, deleteItem, updateItem, isAuthenticated, setIsAuthenticated, userRole }}>
            {children}
        </ItemContext.Provider>
    );
};
