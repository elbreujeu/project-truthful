import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { API_URL, clientId } from '../Env';
import '../styles/style.css';
import '../styles/Login.css';
import { GoogleLogin } from '@react-oauth/google';

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
        const errorMessage = returnMessage.error.charAt(0).toUpperCase() + returnMessage.error.slice(1);
        console.error(returnMessage)
        setError(errorMessage);
      }
    } catch (error) {
      console.error(error);
      setError("An error occurred while logging in, please try again later");
    }
  };

  const handleGoogleSuccess = async (response) => {
    console.log(response);
  };

  const handleGoogleFailure = async (response) => {
    console.log(response);
    setError("An error occurred while logging in with Google, please try again later");
  };

  return (
    <div className="background-color" style={{display: 'flex', flexDirection: 'column', alignItems: 'center', height: '100vh'}}>
      <h1 className="text" style={{marginTop: '3rem'}}>Sign in</h1>
      {error && <div className="error_box" style={{marginBottom: '1rem'}}>{error}</div>}
      <div style={{width: '30%', display: 'flex', flexDirection: 'column', alignItems: 'flex-start'}}>
        <label className='text' style={{alignSelf: 'flex-start', marginBottom: '0.5rem'}}>Username</label>
          <input 
              type="text" 
              id="username"
              className='text_box'
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === 'Enter') {
                  handleLogin();
                }
              }}
          />
        <label className='text' style={{alignSelf: 'flex-start', marginTop:'0.5rem', marginBottom: '0.5rem'}}>Password</label>
        <input 
            type="password" 
            id="password" 
            className='text_box'
            value={password} 
            onChange={(e) => setPassword(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                handleLogin();
              }
            }}
        />
        <button 
            className="button" 
            style={{padding: '1rem', marginTop: '2rem', alignSelf: 'center', fontFamily: 'Fira Code', fontSize: '1rem'}}
            onClick={handleLogin}
        >
            Login
        </button>
        <GoogleLogin
          buttonText="Login with Google"
          clientId={clientId}
          onSuccess={(response) => handleGoogleSuccess(response)}
          onError={(response) => handleGoogleFailure(response)}
          cookiePolicy={'single_host_origin'}
        />
        <p className="alt_text" style={{marginTop: '0rem', alignSelf: 'center'}}>Forgot password ?</p>
      </div>
    </div>
  );
};

export default Login;