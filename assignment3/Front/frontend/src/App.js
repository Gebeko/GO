import React, { useState, useEffect, useContext } from 'react';
import ItemList from './Components/ItemList';
import { ItemProvider} from './Context/ItemContext';
import LoginForm from './Components/LoginForm';
import { jwtDecode } from 'jwt-decode';

const App = () => {
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [userRole, setUserRole] = useState(null);

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (token) {
            setIsAuthenticated(true);
            const decodedToken = jwtDecode(token);
            setUserRole(decodedToken.role);
        }
    }, []);

    const handleLogin = (token) => {
        setIsAuthenticated(true);
        localStorage.setItem('token', token);
        const decodedToken = jwtDecode(token);
        console.log("Decoded Token:", decodedToken);
        setUserRole(decodedToken.role);
    };

    const handleLogout = () => {
        setIsAuthenticated(false);
        setUserRole(null);
        localStorage.removeItem('token');
    };

    return (
        <ItemProvider>
            <div className="App">
                {isAuthenticated ? (
                    <>  
                        <ItemList userRole={userRole} />
                        <button onClick={handleLogout}>Logout</button>
                    </>
                ) : (
                    <LoginForm onLogin={handleLogin} />
                )}
            </div>
        </ItemProvider>
    );
};

export default App;
