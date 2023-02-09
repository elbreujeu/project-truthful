import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { API_URL } from '../Env';
import '../styles/colors.css';

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleLogin = async () => {
    try {
      const response = await fetch(API_URL + '/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      });
      if (response.status === 200) {
        const { token } = await response.json();
        document.cookie = `token=${token}`;
        navigate('/feed');
      } else {
        const returnMessage = await response.json();
        // sets error message as returnMessage.error but capitalizes the first letter
        const errorMessage = returnMessage.error.charAt(0).toUpperCase() + returnMessage.error.slice(1);
        console.error(returnMessage)
        setError(errorMessage);
      }
    } catch (error) {
      console.error(error);
      setError("An error occurred while logging in, please try again later");
    }
};

  return (
    <div className="background-color" style={{display: 'flex', flexDirection: 'column', alignItems: 'center', height: '100vh'}}>
      <h1 className="text" style={{marginTop: '3rem'}}>Sign in</h1>
      {error && <div className="error_box">{error}</div>}
      <div style={{width: '30%', display: 'flex', flexDirection: 'column', alignItems: 'flex-start'}}>
        <label className='text' style={{alignSelf: 'flex-start', marginBottom: '0.5rem'}}>Username</label>
          <input 
              type="text" 
              id="username"
              className='text_box'
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              style={{width: '100%', padding: '0.5rem', marginTop: '0'}} 
          />
        <label className='text' style={{alignSelf: 'flex-start', marginBottom: '0.5rem'}}>Password</label>
        <input 
            type="password" 
            id="password" 
            className='text_box'
            value={password} 
            onChange={(e) => setPassword(e.target.value)}
            style={{width: '100%', padding: '0.5rem', marginTop: '0'}} 
        />
        <button 
            className="button" 
            style={{width: '100%', padding: '1rem', marginTop: '2rem'}}
            onClick={handleLogin}
        >
            Login
        </button>
        <p className="alt_text" style={{marginTop: '0rem'}}>Forgot password ?</p>
      </div>
    </div>
  );
};

export default Login;