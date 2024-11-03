import React, { useState } from 'react';
import axios from 'axios';

const LoginForm = ({ onLogin }) => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [role, setRole] = useState('user');
    const [isRegistering, setIsRegistering] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            if (isRegistering) {
                // Register with role
                await axios.post('http://localhost:8080/register', { username, password, role });
                alert("Registration successful! You can now log in.");
                setIsRegistering(false);
            } else {
                // Log in
                const response = await axios.post('http://localhost:8080/login', { username, password });
                const token = response.data.token;
                localStorage.setItem('token', token);
                onLogin(token);
            }
        } catch (error) {
            console.error("Error:", error);
            alert(isRegistering ? "Registration failed. Please try again." : "Login failed. Check your credentials.");
        }
    };

    return (
        <div>
            <h2>{isRegistering ? "Register" : "Login"}</h2>
            <form onSubmit={handleSubmit}>
                <input
                    type="text"
                    placeholder="Username"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    required
                />
                <input
                    type="password"
                    placeholder="Password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                />
                {isRegistering && (
                    <div>
                        <label>
                            <input
                                type="radio"
                                value="admin"
                                checked={role === 'admin'}
                                onChange={() => setRole('admin')}
                            />
                            Admin
                        </label>
                        <label>
                            <input
                                type="radio"
                                value="user"
                                checked={role === 'user'}
                                onChange={() => setRole('user')}
                            />
                            User
                        </label>
                    </div>
                )}
                <button type="submit">{isRegistering ? "Register" : "Login"}</button>
            </form>
            <button onClick={() => setIsRegistering(!isRegistering)}>
                {isRegistering ? "Already have an account? Log in" : "Don't have an account? Register"}
            </button>
        </div>
    );
};

export default LoginForm;
