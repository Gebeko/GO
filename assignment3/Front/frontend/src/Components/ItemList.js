import React, { useState } from 'react';
import { useItemContext } from '../Context/ItemContext';

const AddItemForm = ({ userRole }) => {
    const [name, setName] = useState('');
    const [price, setPrice] = useState('');
    const { addItem } = useItemContext();

    const handleSubmit = async (e) => {
        e.preventDefault();
        const newItem = { name, price: parseInt(price) };
        await addItem(newItem);
        setName('');
        setPrice('');
    };
    
    return userRole === 'admin' ? (
        <form onSubmit={handleSubmit}>
            <input
                type="text"
                placeholder="Name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
            />
            <input
                type="number"
                placeholder="Price"
                value={price}
                onChange={(e) => setPrice(e.target.value)}
                required
            />
            <button type="submit">Add Item</button>
        </form>
    ) : null;
};

const ItemList = () => {
    const { items, loading, error, deleteItem, updateItem, userRole } = useItemContext();
    const [editItemId, setEditItemId] = useState(null);
    const [editName, setEditName] = useState('');
    const [editPrice, setEditPrice] = useState('');

    const handleEditSubmit = async (e) => {
        e.preventDefault();
        const updatedItem = { id: editItemId, name: editName, price: parseInt(editPrice) };
        await updateItem(updatedItem);
        setEditItemId(null);
        setEditName('');
        setEditPrice('');
    };

    if (loading) return <div>Loading...</div>;
    if (error) return <div>{error}</div>;

    return (
        <div>
            <h2>Items List</h2>
            <AddItemForm userRole={userRole} />
            <ul>
                {items.map(item => (
                    <li key={item.id}>
                        {item.name} - ${item.price}
                        {userRole === 'admin' && (
                            <span>
                                <button onClick={() => {
                                    setEditItemId(item.id);
                                    setEditName(item.name);
                                    setEditPrice(item.price);
                                }}>Edit</button>
                                <button onClick={() => deleteItem(item.id)}>Delete</button>
                            </span>
                        )}
                    </li>
                ))}
            </ul>
            {editItemId && (
                <form onSubmit={handleEditSubmit}>
                    <input
                        type="text"
                        placeholder="Edit Name"
                        value={editName}
                        onChange={(e) => setEditName(e.target.value)}
                        required
                    />
                    <input
                        type="number"
                        placeholder="Edit Price"
                        value={editPrice}
                        onChange={(e) => setEditPrice(e.target.value)}
                        required
                    />
                    <button type="submit">Update Item</button>
                    <button type="button" onClick={() => setEditItemId(null)}>Cancel</button>
                </form>
            )}
        </div>
    );
};

export default ItemList;
